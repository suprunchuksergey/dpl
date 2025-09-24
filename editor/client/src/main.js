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
      <div class="panel w-full h-full flex flex-col">
        <div id="editor" class="h-full w-full"></div>

        <footer
          class="flex justify-end p-2 border-t border-t-blue-100 bg-white"
        >
          <button
            id="run"
            class="bg-blue-600 text-white py-1 px-3 rounded-md cursor-pointer font-medium text-lg"
          >
            Запустить
          </button>
        </footer>
      </div>

      <div class="w-full h-full flex flex-col gap-y-3">
        <div class="panel w-full h-full p-3 flex items-center justify-center">
          <div id="display" class="h-full w-full"></div>
        </div>

        <div
          id="output"
          class="panel w-full h-1/3 shrink-0 overflow-auto whitespace-pre-wrap break-words p-3"
        ></div>
      </div>
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
