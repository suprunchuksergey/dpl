import "./style.css";
import * as monaco from "monaco-editor";
import "/public/wasm_exec.js";

const go = new Go();
WebAssembly.instantiateStreaming(
  fetch("/public/dpl.wasm"),
  go.importObject,
).then((result) => {
  go.run(result.instance);
});

document.querySelector("#app").innerHTML = `
    <div class="h-screen w-screen bg-blue-50 flex gap-x-3 p-3">
      <div
        class="w-full h-full bg-white shadow-sm border border-blue-100 flex flex-col"
      >
        <div id="editor" class="h-full w-full"></div>

        <footer class="flex justify-end p-1 border-t border-t-blue-100">
          <button
            id="run"
            class="bg-blue-600 text-white py-1 px-3 rounded-sm cursor-pointer font-medium"
          >
            Запустить
          </button>
        </footer>
      </div>

      <div
        id="output"
        class="whitespace-pre-wrap break-words w-full h-full overflow-auto bg-[#f5f5f5] shadow-sm border border-blue-100 p-3"
      ></div>
    </div>
`;

monaco.languages.register({ id: "dpl" });

monaco.languages.setLanguageConfiguration("dpl", {
  autoClosingPairs: [
    { open: "(", close: ")" },
    { open: "[", close: "]" },
    { open: "{", close: "}" },
    { open: '"', close: '"' },
  ],
});

monaco.languages.setMonarchTokensProvider("dpl", {
  tokenizer: {
    root: [
      [/\b(if|elif|else|for|in|return|true|false|null)\b/, "keyword"],
      [/\b(and|or|not)\b/, "operator.logical"],
      [/[+\-*\/%]|\|\|/, "operator.arithmetic"],
      [/==|!=|<|>|<=|>=/, "operator.comparison"],
      [/:?=|=/, "operator.assignment"],
      [/[\(\)\[\]\{\}]|;|,|\.|->/, "delimiter"],
      [/[a-zA-Z_][a-zA-Z0-9_]*/, "identifier"],
      [/\d+\.\d*|\.\d+|\d+/, "number"],
      [/"[^"]*"/, "string"],
    ],
  },
});

monaco.editor.defineTheme("dpl", {
  base: "vs",
  inherit: true,
  rules: [
    { token: "keyword", foreground: "#0000ff", fontStyle: "bold" },
    { token: "operator.logical", foreground: "#d35400" },
    { token: "operator.arithmetic", foreground: "#2e7d32" },
    { token: "operator.comparison", foreground: "#2e7d32" },
    { token: "operator.assignment", foreground: "#1b5e20" },
    { token: "delimiter", foreground: "#37474f" },
    { token: "identifier", foreground: "#000000" },
    { token: "number", foreground: "#6a1b9a" },
    { token: "string", foreground: "#c62828" },
  ],
  colors: {
    "editor.background": "#f5f5f5",
    "editor.foreground": "#333333",
  },
});

const editor = monaco.editor.create(document.getElementById("editor"), {
  value: `
factorial := (n) -> {
	if n <= 1 {
		return 1;
	};
	return n * factorial(n-1);
};

sum := 0;

for i in 8 {
	sum = sum + factorial(i);
};

println(sum);
`,
  theme: "dpl",
  language: "dpl",
});

const output = document.getElementById("output");

const run = document.getElementById("run");
run.onclick = () => {
  output.innerText = "";
  exec(editor.getValue(), (data) => {
    output.innerText += data;
  });
};
