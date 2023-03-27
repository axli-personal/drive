import { marked } from "marked";
import hljs from "highlight.js";
import "/src/css/highlight.scss";

marked.setOptions({
  highlight: function (code, lang) {
    const language = hljs.getLanguage(lang) ? lang : "plaintext";
    return hljs.highlight(code, { language }).value;
  },
  langPrefix: "hljs language-"
});

export default marked;