const App = {
  data() {
    return {
      showingEditor: true,
    };
  },
  methods: {
    // TODO: differentiate switching to preview vs generating new preview
    async switchEditor() {
      if (!this.showingEditor) {
        this.showingEditor = true;
        this.previewContent = "";
        return;
      }
      const resp = await fetch("../admin/blog/preview", {
        method: "POST",
        body: editor.getValue(),
      });
      const text = await resp.text();
      if (!resp.ok) {
        alert(`Error response (${resp.statusText}): ${text}`);
        return;
      }
      document.querySelector("#preview").innerHTML = text;
      this.showingEditor = false;
    },
    __blank() {}
  }
};
const app = Vue.createApp(App);
app.mount("#app");
