import "./style.css";
import * as monaco from "monaco-editor";
import "/public/wasm_exec.js";
import { Chart, registerables } from "chart.js";

const go = new Go();
WebAssembly.instantiateStreaming(fetch("/dpl.wasm"), go.importObject).then(
  (result) => {
    go.run(result.instance);
    run.click();
  },
);

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

      <div class="w-full h-full flex flex-col gap-y-3 overflow-hidden">
        <div class="panel w-full aspect-video p-3 flex items-center justify-center">
          <canvas id="display" class="h-full w-full"></canvas>
        </div>

        <div
          id="output"
          class="panel w-full h-full overflow-auto whitespace-pre-wrap break-words p-3"
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
  value: `power = (base, exponent) -> {
    if exponent == 0 { return 1 };
    if exponent == 1 { return base };
    return base * power(base, exponent - 1)
};

println("подготовка данных");

data = [];

for i in 12 {
    data = append(
        data,
        {
            "число": i,
            "треугольные числа": i * (i + 1) / 2,
            "чередование с ростом": i + power(-1, i),
            "пентагональные числа": (3 * power(i, 2) - i) / 2,
        }
    )
};

println("рендеринг данных");

draw(
    "line",
    data,
    {
        "id": "число",
        "values": [
            "треугольные числа",
            "чередование с ростом", 
            "пентагональные числа",
        ],
        "title": "Сравнение математических последовательностей",
    }
);
`,
  theme: "dpl",
  language: "dpl",
  fontWeight: "500",
});

const output = document.getElementById("output");

const write = (data) => {
  output.innerText += data;
};

Chart.defaults.font.family = "firacode";
Chart.defaults.font.size = 14;
Chart.defaults.font.weight = "bold";

const GRID_COLOR = "#d8d8d8";

const COLORS = ["#377ec0", "#f04f52", "#12baaa", "#5460ac", "#f7891f"];
const COLORS2 = ["#aa78ff", "#ffa5eb", "#d7ff3c", "#6ef0a0", "#003c3c"];

const isOneColor = (type) => {
  return ["bar", "line", "radar", "scatter", "bubble"].includes(type);
};

const hasGrid = (type) => {
  return ["bar", "bubble", "line", "scatter"].includes(type);
};

const elByIndex = (arr, index) => {
  return arr[index % arr.length];
};

Chart.register(...registerables);

const display = document.getElementById("display");

let chart = null;

const draw = (type, data, options) => {
  if (chart) chart.destroy();

  chart = new Chart(display, {
    type,
    data: {
      labels: data.map((row) => row[options.id]),
      datasets: options.values.map((v, index) => {
        return {
          label: v,
          data: data.map((row) => row[v]),

          backgroundColor: isOneColor(type)
            ? type === "radar"
              ? elByIndex(COLORS, index) + "75"
              : elByIndex(COLORS, index)
            : elByIndex([COLORS, COLORS2], index),

          borderColor: isOneColor(type)
            ? elByIndex(COLORS, index)
            : elByIndex([COLORS, COLORS2], index),
        };
      }),
    },
    options: {
      aspectRatio: 2,
      plugins: {
        legend: { position: "bottom" },
        title: {
          display: true,
          text: options.title,
          font: { size: 16 },
        },
      },
      scales: {
        x: {
          display: hasGrid(type),
          grid: { color: GRID_COLOR },
        },
        y: {
          display: hasGrid(type),
          grid: { color: GRID_COLOR },
        },
      },
    },
  });
};

const run = document.getElementById("run");

run.onclick = () => {
  output.innerText = "";

  exec(editor.getValue(), write, draw);
};
