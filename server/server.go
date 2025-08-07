package server

import (
	"context"
	"crypto/tls"
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

	"github.com/caddyserver/certmagic"
	"github.com/johnietre/my-site/server/blogs"
	"github.com/johnietre/my-site/server/handlers"
	"github.com/johnietre/my-site/server/products"
	"github.com/johnietre/my-site/server/repos"
	utils "github.com/johnietre/utils/go"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	zapLogger *zap.Logger
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
		"products-db", "$(base-dir)/products/products.db", "Path to products database",
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

	flags.Bool("no-repos", false, "Disable repository refreshing")
	flags.Bool(
		"autotls",
		false,
		"TLS automation through certmagic (provide domain name rather than IP:PORT address); the email address used to for the ACME server account can be specified with the MY_SITE_ACME_EMAIL environment variable",
	)
	cmd.MarkFlagsMutuallyExclusive("autotls", "cert")
	return cmd
}

func run(cmd *cobra.Command, args []string) {
	flags := cmd.Flags()

	addr := flags.Arg(0)
	if addr == "" {
		addr = "127.0.0.1:8000"
	}

	baseDir := utils.Must(flags.GetString("base-dir"))

	certPath := utils.Must(flags.GetString("cert"))
	keyPath := utils.Must(flags.GetString("key"))

	staticDir := utils.Must(flags.GetString("static-dir"))
	blogsDir := utils.Must(flags.GetString("blogs-dir"))
	tmplsDir := utils.Must(flags.GetString("tmpls-dir"))

	productsDbPath := utils.Must(flags.GetString("products-db"))
	blogsDbPath := utils.Must(flags.GetString("blogs-db"))

	adminConfigPath := utils.Must(flags.GetString("admin-config"))
	remoteIP := utils.Must(flags.GetString("remote-ip"))

	logPath := utils.Must(flags.GetString("log"))

	noRepos := utils.Must(flags.GetBool("no-repos"))

	autotls := utils.Must(flags.GetBool("autotls"))

	if logPath != "" {
		f, err := utils.OpenAppend(logPath)
		if err != nil {
			log.Fatalf("error opening log file: %v", err)
		}
		w := zapcore.Lock(f)
		log.SetOutput(w)
		zapLogger = zap.New(zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig()),
			w,
			zap.InfoLevel,
		))
	}

	if baseDir == "" {
		baseDir = "./"
	}
	baseDir = filepath.Join(baseDir, "a")
	// Add a dummy part to the path to add a trailing slug
	baseDir = baseDir[:len(baseDir)-1]

	replacePath(&adminConfigPath, "$(base-dir)/", baseDir)
	replacePath(&staticDir, "$(base-dir)/", baseDir)
	replacePath(&productsDbPath, "$(base-dir)/", baseDir)
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
	} else if err = products.InitProducts(productsDbPath); err != nil {
		log.Fatalf("error initializing products: %v", err)
	} else if err = blogs.InitBlogs(blogsDir, blogsDbPath); err != nil {
		log.Fatalf("error initializing blogs: %v", err)
	} else if err = repos.InitRepos(!noRepos); err != nil {
		log.Fatalf("error initializing blogs: %v", err)
	}

	doneChan := make(chan bool, 2)
	err = RunServer(RunServerConfig{
		Addr:      addr,
		StaticDir: staticDir,
		KeyPath:   keyPath,
		CertPath:  certPath,
		Autotls:   autotls,
		ACMEEmail: os.Getenv("MY_SITE_ACME_EMAIL"),
	}, doneChan)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("error running server: %v", err)
	}
	<-doneChan
	log.Print("finished")
}

func RunServer(
	cfg RunServerConfig,
	doneChan chan<- bool,
) error {
	// Most of this code was copied from certmagic.HTTPS/certmagic.TLS functions.
	ctx := context.Background()

	var httpLn net.Listener
	var tlsConf *tls.Config
	var certCfg *certmagic.Config
	if cfg.Autotls {
		certmagic.DefaultACME.Agreed = true
		certmagic.DefaultACME.Email = cfg.ACMEEmail
		if zapLogger != nil {
			certmagic.DefaultACME.Logger = zapLogger
			certmagic.Default.Logger = zapLogger
		}
		certCfg = certmagic.NewDefault()
		tlsConf = certCfg.TLSConfig()
		if err := certCfg.ManageSync(ctx, []string{cfg.Addr}); err != nil {
			return err
		}
		lc, err := newListenConfig(), error(nil)
		httpLn, err = lc.Listen(ctx, "tcp", ":80")
		if err != nil {
			log.Fatalf("error starting listener: %v", err)
		}
		tlsConf.NextProtos = append(
			[]string{"h2", "http/1.1"},
			tlsConf.NextProtos...,
		)
		cfg.Addr = ":443"
	}

	lc := newListenConfig()
	httpsLn, err := lc.Listen(ctx, "tcp", cfg.Addr)
	if err != nil {
		log.Fatalf("error starting listener: %v", err)
	}
	if cfg.Autotls {
		httpsLn = tls.NewListener(httpsLn, tlsConf)
	}

	var httpSrvr, httpsSrvr *http.Server
	if httpLn != nil {
		httpSrvr = &http.Server{
			ReadHeaderTimeout: 10 * time.Second,
			ReadTimeout:       10 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       10 * time.Second,
			BaseContext:       func(ln net.Listener) context.Context { return ctx },
		}
		if certCfg != nil && len(certCfg.Issuers) != 0 {
			if am, ok := certCfg.Issuers[0].(*certmagic.ACMEIssuer); ok {
				httpSrvr.Handler = am.HTTPChallengeHandler(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						reqHost, _, err := net.SplitHostPort(r.Host)
						if err != nil {
							reqHost = r.Host
						}
						to := reqHost + r.URL.RequestURI()
						w.Header().Set("Connection", "close")
						http.Redirect(w, r, to, http.StatusMovedPermanently)
					},
				))
			}
		}
	}
	if httpsLn != nil {
		httpsSrvr = &http.Server{
			Handler:           handlers.CreateRouter(cfg.StaticDir),
			ReadHeaderTimeout: 10 * time.Second,
			ReadTimeout:       10 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       10 * time.Second,
			BaseContext:       func(ln net.Listener) context.Context { return ctx },
			TLSConfig:         tlsConf,
		}
	}

	interruptChan := make(chan os.Signal, 5)
	go func() {
		<-interruptChan
		log.Print("shutting down")
		go func() {
			if httpSrvr != nil {
				httpSrvr.Shutdown(context.Background())
			}
			if httpsSrvr != nil {
				httpsSrvr.Shutdown(context.Background())
			}
			doneChan <- true
		}()
		<-interruptChan
		log.Print("closing")
		if httpSrvr != nil {
			httpSrvr.Close()
		}
		if httpsSrvr != nil {
			httpsSrvr.Close()
		}
		doneChan <- true
	}()
	signal.Notify(interruptChan, os.Interrupt)

	if httpsLn != nil {
		log.Printf("running on %s", httpsLn.Addr())
	}
	if httpLn != nil {
		log.Printf("running HTTP on %s", httpLn.Addr())
	}

	if cfg.CertPath != "" || cfg.KeyPath != "" {
		err = httpsSrvr.ServeTLS(httpsLn, cfg.CertPath, cfg.KeyPath)
	} else {
		if httpsSrvr != nil {
			err = httpsSrvr.Serve(httpsLn)
		}
		if err == nil && httpSrvr != nil {
			err = httpSrvr.Serve(httpLn)
		}
	}
	return err
}

func newListenConfig() net.ListenConfig {
	return net.ListenConfig{
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

type RunServerConfig struct {
	Addr                         string
	StaticDir, KeyPath, CertPath string
	Autotls                      bool
	ACMEEmail                    string
}

func replacePath(strPtr *string, from, to string) {
	if strings.HasPrefix(*strPtr, from) {
		*strPtr = filepath.Clean(strings.Replace(*strPtr, from, to, 1))
	}
}
