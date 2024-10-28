package server

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/johnietre/my-site/server/apps"
	"github.com/johnietre/my-site/server/blogs"
	"github.com/johnietre/my-site/server/handlers"
	"github.com/johnietre/my-site/server/repos"
	utils "github.com/johnietre/utils/go"
	"github.com/spf13/cobra"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
}

func Run() {
	rootCmd := &cobra.Command{
		Use:   "my-site",
		Short: "My website",
		Long:  "My website.",
	}
	rootCmd.AddCommand(makeRunCmd(), makeConfigCmd())
	if err := rootCmd.Execute(); err != nil {
		log.SetFlags(0)
		log.Fatal(err)
	}
}

func makeRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "run [FLAGS] [ADDR (default: 127.0.0.1:8000)]",
		Short:                 "Run the server",
		Long:                  "Run the server.",
		Args:                  cobra.MaximumNArgs(1),
		DisableFlagsInUseLine: true,
		Run:                   run,
	}
	flags := cmd.Flags()

	//flags.String("addr", "127.0.0.1:8000", "Address to run on")
	flags.String("cert", "", "Path to cert file")
	flags.String("key", "", "Path to key file")

	flags.String(
		"base-dir", ".",
		"Path to base directory where all other directories are located",
	)

	flags.String("static-dir", "$(base-dir)/static", "Path to static directory")
	flags.String(
		"apps-db", "$(base-dir)/apps/apps.db", "Path to apps database",
	)
	flags.String(
		"blogs-dir", "$(base-dir)/blogs", "Path to blogs directory",
	)
	flags.String(
		"blogs-db", "$(base-dir)/blogs/blogs.db", "Path to blogs database",
	)
	flags.String(
		"tmpls-dir", "$(base-dir)/templates", "Path to templates directory",
	)

	flags.String(
		"admin-config",
		"$(base-dir)/server/admin-config.json",
		"Path to admin config file",
	)
	flags.String("remote-ip", "", "Remote IP to check for to parse")

	flags.String("log", "", "Path to log file (empty routes to stderr)")
	return cmd
}

func run(cmd *cobra.Command, args []string) {
	flags := cmd.Flags()

	addr := flags.Arg(0)
	if addr == "" {
		addr = "127.0.0.1:8000"
	}

	baseDir, _ := flags.GetString("base-dir")

	certPath, _ := flags.GetString("cert")
	keyPath, _ := flags.GetString("key")

	staticDir, _ := flags.GetString("static-dir")
	blogsDir, _ := flags.GetString("blogs-dir")
	tmplsDir, _ := flags.GetString("tmpls-dir")

	appsDbPath, _ := flags.GetString("apps-db")
	blogsDbPath, _ := flags.GetString("blogs-db")

	adminConfigPath, _ := flags.GetString("admin-config")
	remoteIP, _ := flags.GetString("remote-ip")

	logPath, _ := flags.GetString("log")

	if logPath != "" {
		f, err := utils.OpenAppend(logPath)
		if err != nil {
			log.Fatalf("error opening log file: %v", err)
		}
		log.SetOutput(f)
	}

	if baseDir == "" {
		baseDir = "./"
	}
	baseDir = filepath.Join(baseDir, "a")
	// Add a dummy part to the path to add a trailing slug
	baseDir = baseDir[:len(baseDir)-1]

	replacePath(&adminConfigPath, "$(base-dir)/", baseDir)
	replacePath(&staticDir, "$(base-dir)/", baseDir)
	replacePath(&appsDbPath, "$(base-dir)/", baseDir)
	replacePath(&blogsDir, "$(base-dir)/", baseDir)
	replacePath(&blogsDbPath, "$(base-dir)/", baseDir)
	replacePath(&tmplsDir, "$(base-dir)/", baseDir)

	adminConfig := handlers.AdminConfig{}
	f, err := os.Open(adminConfigPath)
	if err != nil {
		log.Fatalf("error opening admin config file: %v", err)
	}
	err = json.NewDecoder(f).Decode(&adminConfig)
	f.Close()
	if err != nil {
		log.Fatalf("error decoding admin config file: %v", err)
	}

	if err := handlers.InitHandlers(tmplsDir, remoteIP, adminConfig); err != nil {
		log.Fatalf("error initializing handlers: %v", err)
	} else if err = apps.InitApps(appsDbPath); err != nil {
		log.Fatalf("error initializing apps: %v", err)
	} else if err = blogs.InitBlogs(blogsDir, blogsDbPath); err != nil {
		log.Fatalf("error initializing blogs: %v", err)
	} else if err = repos.InitRepos(); err != nil {
		log.Fatalf("error initializing blogs: %v", err)
	}

	doneChan := make(chan bool, 2)
	err = RunServer(addr, staticDir, keyPath, certPath, doneChan)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("error running server: %v", err)
	}
	<-doneChan
	log.Print("finished")
}

func RunServer(addr, staticDir, keyPath, certPath string, doneChan chan<- bool) error {
	lc := net.ListenConfig{
		Control: func(network, addr string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				if err := syscall.SetsockoptInt(
					int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1,
				); err != nil {
					log.Fatalf("error setting SO_REUSEADDR: %v", err)
				}
			})
		},
	}
	ls, err := lc.Listen(context.Background(), "tcp", addr)
	if err != nil {
		log.Fatalf("error starting listener: %v", err)
	}

	srvr := &http.Server{
		Addr:              addr,
		Handler:           handlers.CreateRouter(staticDir),
		ReadHeaderTimeout: time.Second * 10,
	}

	interruptChan := make(chan os.Signal, 5)
	go func() {
		<-interruptChan
		log.Print("shutting down")
		go func() {
			srvr.Shutdown(context.Background())
			doneChan <- true
		}()
		<-interruptChan
		log.Print("closing")
		srvr.Close()
		doneChan <- true
	}()
	signal.Notify(interruptChan, os.Interrupt)

	log.Printf("Running on %s", srvr.Addr)
	if certPath != "" || keyPath != "" {
		err = srvr.ServeTLS(ls, certPath, keyPath)
	} else {
		err = srvr.Serve(ls)
	}
	return err
}

func makeConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "config",
		Short:                 "Config related stuff",
		Long:                  "Config related stuff.",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			log.SetFlags(0)

			flags := cmd.Flags()

			if path, _ := flags.GetString("new-admin"); path != "" {
				if info, err := os.Stat(path); err != nil {
					if !os.IsNotExist(err) {
						log.Fatal(err)
					}
					path = filepath.Join(path, "admin-config.json")
				} else if info.IsDir() {
					path = filepath.Join(path, "admin-config.json")
				}
				f, err := os.Create(path)
				if err != nil {
					log.Fatalf("error creating admin config file: %v", err)
				}
				defer f.Close()
				e := json.NewEncoder(f)
				e.SetIndent("", "  ")
				if err := e.Encode(handlers.AdminConfig{}); err != nil {
					log.Fatalf("error creating admin config file: %v", err)
				}
			}
		},
	}
	flags := cmd.Flags()

	flags.String("new-admin", "", "Create admin config file")
	return cmd
}

func replacePath(strPtr *string, from, to string) {
	if strings.HasPrefix(*strPtr, from) {
		*strPtr = filepath.Clean(strings.Replace(*strPtr, from, to, 1))
	}
}
