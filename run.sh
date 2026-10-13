#!/usr/bin/env bash
# Roda cada versão N vezes e salva as saídas e um resumo em CSV.
# Uso: ./run.sh [execuções]   (padrão: 5)
set -euo pipefail
cd "$(dirname "$0")"

RUNS=${1:-5}

# Lê uma constante de tempo do config (ex.: "20 * time.Millisecond") em ms
const_ms() {
    awk -v name="$1" '$1 == name {
        n = $3
        if ($5 ~ /Second$/ && $5 !~ /Milli/) n *= 1000
        print n
    }' config/config.go
}
const_int() { awk -v name="$1" '$1 == name { print $3 }' config/config.go; }

# Uma pasta por configuração, para uma bateria não apagar a outra
OUT="results/think$(const_ms ThinkTime)_eat$(const_ms EatTime)_n$(const_int Philosophers)_r$(const_int Rounds)"
echo "Salvando em $OUT"
rm -rf "$OUT"
mkdir -p "$OUT" bin

# Ambiente de execução (pedido no relatório)
{
    go version
    grep PRETTY_NAME /etc/os-release
    lscpu | grep -E "Model name|^CPU\(s\):|Core\(s\) per socket|Thread\(s\) per core"
    grep -E "Time|Philosophers|Rounds" config/config.go
} > "$OUT/env.txt"

for v in v0 v1 v2 v3a v3b; do
    go build -o "bin/$v" "./$v"
done

# V1: uma execução só, como evidência (pode travar; o timeout evita esperar para sempre)
timeout 30 "bin/v1" > "$OUT/v1.txt" 2>&1 || true

# Converte as métricas de uma saída em uma linha CSV (durações em ms)
summarize() {
    awk -v ver="$1" -v run="$2" '
        function ms(d,    a) {
            if (d ~ /^[0-9]+m[0-9]/) { split(d, a, "m"); return a[1] * 60000 + ms(a[2]) }
            if (d ~ /ns$/) return d / 1e6
            if (d ~ /µs$/) return d / 1e3
            if (d ~ /ms$/) return d + 0
            if (d ~ /s$/)  return d * 1000
            return d + 0
        }
        function list(s,    a, n, i, r) {
            gsub(/[\[\]]/, "", s)
            n = split(s, a, " ")
            for (i = 1; i <= n; i++) r = r (i > 1 ? " " : "") sprintf("%.3f", ms(a[i]))
            return r
        }
        /took/              { total = ms($NF) }
        /^Mean Wait Time:/  { mean = ms($NF) }
        /^Max Wait Time:/   { max = ms($NF) }
        /^Max Wait Per/     { sub(/.*: /, ""); per = list($0) }
        /^Cycles per/       { sub(/.*: /, ""); gsub(/[\[\]]/, ""); cycles = $0 }
        /^Max eating/       { eating = $NF }
        END { printf "%s,%s,%.3f,%.3f,%.3f,%s,%s,%s\n", ver, run, total, mean, max, eating, cycles, per }
    ' "$3"
}

CSV="$OUT/summary.csv"
echo "version,run,total_ms,mean_wait_ms,max_wait_ms,max_eating,cycles,max_wait_per_philosopher_ms" > "$CSV"

for v in v0 v2 v3a v3b; do
    for i in $(seq 1 "$RUNS"); do
        echo "$v: execução $i/$RUNS"
        "bin/$v" > "$OUT/${v}_run$i.txt"
        summarize "$v" "$i" "$OUT/${v}_run$i.txt" >> "$CSV"
    done
done

echo "Pronto: $CSV"
