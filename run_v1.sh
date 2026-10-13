#!/usr/bin/env bash
# Roda a V1 N vezes para ver se o deadlock acontece naturalmente.
# Uso: ./run_v1.sh [execuções]   (padrão: 50)
set -euo pipefail
cd "$(dirname "$0")"

RUNS=${1:-50}
OUT=results/v1
mkdir -p "$OUT" bin

# Lê uma constante de tempo do config (ex.: "20 * time.Millisecond") em ms
const_ms() {
    awk -v name="$1" '$1 == name {
        n = $3
        if ($5 ~ /Second$/ && $5 !~ /Milli/) n *= 1000
        print n
    }' config/config.go
}
const_int() { awk -v name="$1" '$1 == name { print $3 }' config/config.go; }

THINK=$(const_ms ThinkTime)
EAT=$(const_ms EatTime)
BETWEEN=$(const_ms TimeBetweenForks)
N=$(const_int Philosophers)
ROUNDS=$(const_int Rounds)

# Pior caso sem deadlock: tudo sequencial, como a V0. Folga de 2x.
EXPECTED=$(( ROUNDS * N * (THINK + EAT) ))
TIMEOUT_S=$(( (EXPECTED * 2 + 999) / 1000 ))

echo "think=${THINK}ms eat=${EAT}ms filósofos=$N rodadas=$ROUNDS entre_garfos=${BETWEEN}ms"
echo "pior caso esperado: ${EXPECTED}ms, timeout: ${TIMEOUT_S}s"
[ "$BETWEEN" != "0" ] && echo "AVISO: TimeBetweenForks != 0, o deadlock está sendo forçado"

go build -o bin/v1 ./v1

CSV="$OUT/runs.csv"
echo "run,result,total_ms,max_eating,meals" > "$CSV"
ok=0; dead=0; hang=0

for i in $(seq 1 "$RUNS"); do
    file="$OUT/run$i.txt"
    set +e
    timeout "$TIMEOUT_S" bin/v1 > "$file" 2>&1
    code=$?
    set -e

    meals=$(grep -cE ' E[0-9]+$' "$file" || true)
    total=""; eating=""
    if grep -q "all goroutines are asleep" "$file"; then
        result=deadlock; dead=$((dead + 1))
    elif [ "$code" = 124 ]; then
        result=timeout; hang=$((hang + 1))
    elif [ "$code" = 0 ]; then
        result=ok; ok=$((ok + 1))
        total=$(awk '/took/ { d = $NF
            if (d ~ /ms$/) print d + 0; else if (d ~ /s$/) print d * 1000 }' "$file")
        eating=$(awk '/^Max eating/ { print $NF }' "$file")
        rm "$file"   # só guarda as saídas interessantes
    else
        result="erro($code)"
    fi

    echo "$i,$result,$total,$eating,$meals" >> "$CSV"
    echo "execução $i/$RUNS: $result (refeições: $meals)"
done

echo
echo "ok: $ok   deadlock: $dead   timeout: $hang   (de $RUNS)"
echo "Detalhes: $CSV   Saídas dos deadlocks/timeouts: $OUT/"
