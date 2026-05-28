import os
import json
import geo
from pathlib import Path

from mpi4py import MPI
from petsc4py import PETSc
from petsc4py.PETSc import ScalarType

import numpy as np
import ufl

from dolfinx import fem, mesh, plot
from dolfinx.fem.petsc import LinearProblem
from dolfinx.io import gmsh as gmshio
from dolfinx import geometry

import pyvista
import matplotlib.pyplot as plt

iid = os.getenv("IID")
config_dir = os.getenv("CONFIG_DIR")
blob_dir = os.getenv("BLOB_DIR")

print(f"IID: {iid}")

with open(f"{config_dir}/sol/mf/{iid}/config.json", "r") as file:
    conf = json.load(file)
    yoke = conf["yoke"]
    coils = conf["coils"]

Path(f"{blob_dir}/sol/mf/{iid}").mkdir(parents=True, exist_ok=True)

if not Path(f"{blob_dir}/sol/mf/{iid}/mesh.msh").is_file():
    geo.create_mesh(f"{blob_dir}/sol/mf/{iid}")

data = gmshio.read_from_msh(
    f"{blob_dir}/sol/mf/{iid}/mesh.msh", MPI.COMM_WORLD, gdim=2)

V = fem.functionspace(data.mesh, ("CG", 1))
Q = fem.functionspace(data.mesh, ("DG", 0))
x = ufl.SpatialCoordinate(data.mesh)
r = x[1]

axis = mesh.locate_entities_boundary(
    data.mesh, dim=1, marker=lambda x: np.isclose(x[1], 0.0))
remote = mesh.locate_entities_boundary(data.mesh, dim=1, marker=lambda x: (
    np.isclose(x[1], 240e-3) | np.isclose(x[0], -
                                          1000e-3) | np.isclose(x[0], 200e-3)
))

axis_dofs = fem.locate_dofs_topological(V, entity_dim=1, entities=axis)
remote_dofs = fem.locate_dofs_topological(V, entity_dim=1, entities=remote)

bc_axis = fem.dirichletbc(value=ScalarType(0), dofs=axis_dofs, V=V)
bc_remote = fem.dirichletbc(value=ScalarType(0), dofs=remote_dofs, V=V)

air = data.cell_tags.find(1)
iron = data.cell_tags.find(2)
nbsn = data.cell_tags.find(3)
nbti = data.cell_tags.find(4)
comp = data.cell_tags.find(5)

nbsn_n = coils["nbsn"]["numberOfTurns"]
nbsn_i = coils["nbsn"]["currentAmp"]
nbsn_s = ((55 - 30) * 1e-3) * ((749.7 - 496.7) * 1e-3)

nbti_n = coils["nbti"]["numberOfTurns"]
nbti_i = coils["nbti"]["currentAmp"]
nbti_s = ((91.5 - 58) * 1e-3) * ((487.7 - 234.7) * 1e-3)

comp_n = coils["comp"]["numberOfTurns"]
comp_i = coils["comp"]["currentAmp"]
comp_s = ((52.4 - 30) * 1e-3) * ((164.7 - 36.7) * 1e-3)

J = fem.Function(Q)
J.x.array[:] = 0.0
J.x.array[nbsn] = fem.Constant(data.mesh, (nbsn_n * nbsn_i) / nbsn_s)
J.x.array[nbti] = fem.Constant(data.mesh, (nbti_n * nbti_i) / nbti_s)
J.x.array[comp] = fem.Constant(data.mesh, (comp_n * comp_i) / comp_s)

mu0 = 4 * np.pi * 1e-7

mu = fem.Function(Q)
mu.x.array[air] = fem.Constant(data.mesh, mu0)
mu.x.array[iron] = fem.Constant(data.mesh, mu0 * yoke["relativePermeability"])
mu.x.array[nbsn] = fem.Constant(data.mesh, mu0)
mu.x.array[nbti] = fem.Constant(data.mesh, mu0)
mu.x.array[comp] = fem.Constant(data.mesh, mu0)

u = ufl.TrialFunction(V)
v = ufl.TestFunction(V)
a = (1.0 / (mu*r)) * ufl.dot(ufl.grad(u), ufl.grad(v)) * ufl.dx
L = J * v * ufl.dx

bcs = [bc_axis, bc_remote]
problem = LinearProblem(
    a,
    L,
    bcs=bcs,
    petsc_options_prefix="solenoid_",
    petsc_options={"ksp_type": "preonly", "pc_type": "lu",
                   "ksp_error_if_not_converged": True},
)

uh = problem.solve()

Bz = fem.Function(Q)
Bz.interpolate(fem.Expression((1.0 / r) * ufl.grad(uh)
               [1], Q.element.interpolation_points))

Br = fem.Function(Q)
Br.interpolate(fem.Expression(-(1.0 / r) * ufl.grad(uh)
               [0], Q.element.interpolation_points))

dz = ufl.Measure("ds", domain=data.mesh, subdomain_data=data.facet_tags)
I = fem.assemble_scalar(fem.form(Bz * dz(8)))

print(I)

plotter = pyvista.Plotter(off_screen=True)

grid = pyvista.UnstructuredGrid(*plot.vtk_mesh(data.mesh))
grid.point_data["u"] = uh.x.array.real
grid.set_active_scalars("u")

plotter.add_mesh(
    grid,
    show_edges=True,
    edge_opacity=0.5,
    show_scalar_bar=False
)

plotter.add_mesh(
    pyvista.UnstructuredGrid(
        *plot.vtk_mesh(data.mesh, 1, data.facet_tags.find(6))),
    color="black",
    line_width=3
)

plotter.add_mesh(
    pyvista.UnstructuredGrid(
        *plot.vtk_mesh(data.mesh, 1, data.facet_tags.find(7))),
    color="black",
    line_width=2
)

plotter.camera_position = "xy"
plotter.reset_camera(bounds=[-1, 0.2, 0, 0.2, 0, 0])
plotter.camera.zoom("tight")
plotter.screenshot(f"{blob_dir}/sol/mf/{iid}/u.png")

plotter = pyvista.Plotter(off_screen=True)

grid = pyvista.UnstructuredGrid(*plot.vtk_mesh(data.mesh))
grid.cell_data["Bz"] = Bz.x.array.real
grid.set_active_scalars("Bz")
grid = grid.cell_data_to_point_data()

plotter.add_mesh(
    grid,
    show_edges=True,
    edge_opacity=0.5,
    show_scalar_bar=False,
    cmap="plasma"
)

plotter.add_mesh(
    pyvista.UnstructuredGrid(*plot.vtk_mesh(
        data.mesh, 1,
        np.concatenate((
            data.facet_tags.find(6),
            data.facet_tags.find(7)
        ))
    )),
    color="black",
    line_width=2
)

plotter.camera_position = "xy"
plotter.reset_camera(bounds=[-1, 0.2, 0, 0.2, 0, 0])
plotter.camera.zoom("tight")
plotter.screenshot(f"{blob_dir}/sol/mf/{iid}/b.png")

eps = 1e-6
nz = 400
z = np.linspace(-1.0, 0.2, nz)
points = np.zeros((3, nz), dtype=np.float64)
points[0, :] = z
points[1, :] = eps
points[2, :] = 0.0

tree = geometry.bb_tree(data.mesh, 2)
cell_candidates = geometry.compute_collisions_points(tree, points.T)
colliding_cells = geometry.compute_colliding_cells(
    data.mesh, cell_candidates, points.T)

points_on_proc = []
cells = []
z_plot = []

for i, p in enumerate(points.T):
    links = colliding_cells.links(i)
    if len(links) > 0:
        points_on_proc.append(p)
        cells.append(links[0])
        z_plot.append(z[i])

points_on_proc = np.array(points_on_proc, dtype=np.float64)
z_plot = np.array(z_plot, dtype=np.float64)
Bz_vals = Bz.eval(points_on_proc, cells)
Bz_vals = np.array(Bz_vals).reshape(-1)

plt.figure(figsize=(8, 4))
plt.plot(z_plot, Bz_vals, lw=2)
plt.xlabel("z [m]")
plt.ylabel("Bz [T]")
plt.grid(True)
plt.tight_layout()
plt.savefig(f"{blob_dir}/sol/mf/{iid}/btz.png")

R = fem.petsc.assemble_vector(fem.form(
    (1.0 / (mu * r)) * ufl.dot(ufl.grad(uh), ufl.grad(v)) * ufl.dx
    - J * v * ufl.dx
))

fem.petsc.apply_lifting(R, [fem.form(a)], [bcs], x0=[
                        uh.x.petsc_vec], alpha=-1.0)
R.ghostUpdate(addv=PETSc.InsertMode.ADD, mode=PETSc.ScatterMode.REVERSE)
fem.petsc.set_bc(R, bcs, uh.x.petsc_vec, -1.0)

print("Weak Residual Norm: ", R.norm())

W1 = fem.assemble_scalar(
    fem.form(2 * ufl.pi * r * (1.0 / (2 * mu) * (Br**2 + Bz**2)) * ufl.dx))
print("W1: ", W1)

W2 = fem.assemble_scalar(fem.form(ufl.pi * J * uh * ufl.dx))
print("W2: ", W2)
print(abs(W2 - W1))
