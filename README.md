# T1 FPPD: Filósofos Jantando em Go

Todos os comandos rodam a partir da raiz do repositório. Os tempos, o número de filósofos e de rodadas ficam em `config/config.go`.

## Versões

```bash
go run ./v0    # sequencial
go run ./v1    # concorrente com deadlock
go run ./v2    # garçom global (sincronização excessiva)
go run ./v3a   # pares/ímpares
go run ./v3b   # semáforo N-1
```

Para forçar o deadlock da V1, use `TimeBetweenForks = 1 * time.Millisecond` no `config/config.go`. Com `0`, o deadlock é raro.

Checagens:

```bash
go vet ./...
go run -race ./v1    # precisa de gcc (sudo apt install build-essential)
```

## Scripts

```bash
./run.sh 100      # roda v0, v2, v3a e v3b 100x cada (padrão: 5)
./run_v1.sh 500   # roda a v1 500x e conta deadlocks (padrão: 50)
```

- `run.sh` salva em `results/think<ms>_eat<ms>_n<N>_r<rodadas>/`: `summary.csv` (uma linha por execução, tempos em ms), `env.txt` (ambiente e config), a saída de cada execução e uma execução da v1.
- `run_v1.sh` salva em `results/v1/`: `runs.csv` e a saída das execuções que travaram.

## Viewer

```bash
python3 -m http.server
```

Abrir http://localhost:8000/viewer/ e escolher a bateria no seletor. Também dá para abrir `viewer/index.html` direto e arrastar o `summary.csv` e o `env.txt` para a página.
