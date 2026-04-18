FROM conda/miniconda3:latest

WORKDIR /svc
COPY environment.yaml .

RUN conda env create -f environment.yaml

COPY . .
ENTRYPOINT ["conda", "run", "--no-capture-output", "-n", "sol.mf", "python", "src/main.py"]
