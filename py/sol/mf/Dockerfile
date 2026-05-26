FROM condaforge/miniforge3:latest

WORKDIR /svc
COPY environment.yaml .

RUN mamba env create -f environment.yaml
RUN mamba clean -afy

COPY . .
ENTRYPOINT ["mamba", "run", "-n", "sol.mf", "python", "src/main.py"]
