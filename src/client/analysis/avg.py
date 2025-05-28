import os
import re

# Percorso alla cartella contenente i file di log
log_dir = "./logs"  # Cambia questo percorso se necessario

# Espressione regolare per trovare i tempi con unità (ms o s)
pattern = re.compile(r"edited remote document in (\d+(?:\.\d+)?)(ms|s)")

# Lista per raccogliere i tempi convertiti in millisecondi
times_ms = []

# Scorri tutti i file nella cartella
for filename in os.listdir(log_dir):
    if filename.endswith(".log"):
        file_path = os.path.join(log_dir, filename)
        with open(file_path, 'r', encoding='utf-8') as file:
            for line in file:
                match = pattern.search(line)
                if match:
                    value = float(match.group(1))
                    unit = match.group(2)
                    if unit == "s":
                        value *= 1000  # Conversione in millisecondi
                    times_ms.append(value)

# Risultati
num_instances = len(times_ms)
average_time_ms = sum(times_ms) / num_instances if num_instances > 0 else 0

print(f"Numero totale di istanze: {num_instances}")
print(f"Tempo medio (in ms): {average_time_ms:.2f}")
