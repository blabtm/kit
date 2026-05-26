import gmsh
import numpy as np

from mpi4py import MPI
from dolfinx.io import gmsh as gmshio
from dolfinx import plot

import pyvista

yoke_h = (194 - 30) * 1e-3
yoke_w = 828.7e-3
yoke_oy = 30e-3
yoke_main_cut_h = (100 - 30) * 1e-3
yoke_main_cut_w = (759.7 - 224.7) * 1e-3
yoke_comp_cut_h = (73 - 30) * 1e-3
yoke_comp_cut_w = (176.2 - 32) * 1e-3
nbti_coil_h = (91.5 - 58) * 1e-3
nbti_coil_w = (749.7 - 496.7) * 1e-3
nbsn_coil_h = (55 - 30) * 1e-3
nbsn_coil_w = (749.7 - 496.7) * 1e-3
comp_coil_h = (52.4 - 30) * 1e-3
comp_coil_w = (164.7 - 36.7) * 1e-3
l_yoke_corner_w = (828.7 - 792.7) * 1e-3
l_yoke_corner_h = (194 - 155.2) * 1e-3
r_yoke_corner_w = 262e-3
r_yoke_corner_h = (194 - 112.7) * 1e-3


def create_mesh(mesh_dir):
    gmsh.initialize()
    model = gmsh.model()
    model.add("Solenoid")

    yoke = model.occ.add_rectangle(-yoke_w, 0, 0, yoke_w, yoke_h)

    # Cut the coils
    cut = model.occ.add_rectangle(-759.7e-3, 0, 0,
                                  yoke_main_cut_w, yoke_main_cut_h)
    yoke = model.occ.cut([(2, yoke)], [(2, cut)])[0][0][1]
    cut = model.occ.add_rectangle(-176.2e-3, 0, 0,
                                  yoke_comp_cut_w, yoke_comp_cut_h)
    yoke = model.occ.cut([(2, yoke)], [(2, cut)])[0][0][1]

    # Cut the left corner
    p1 = model.occ.add_point(-yoke_w, yoke_h, 0)
    p2 = model.occ.add_point(-yoke_w, yoke_h - l_yoke_corner_h, 0)
    p3 = model.occ.add_point(-yoke_w + l_yoke_corner_w, yoke_h, 0)
    cut = model.occ.add_plane_surface([
        model.occ.add_curve_loop([
            model.occ.add_line(p1, p2),
            model.occ.add_line(p1, p3),
            model.occ.add_line(p2, p3)
        ])
    ])
    yoke = model.occ.cut([(2, yoke)], [(2, cut)])[0][0][1]

    # Cut the right corner
    p1 = model.occ.add_point(0, yoke_h, 0)
    p2 = model.occ.add_point(-r_yoke_corner_w, yoke_h, 0)
    p3 = model.occ.add_point(0, yoke_h - r_yoke_corner_h, 0)
    cut = model.occ.add_plane_surface([
        model.occ.add_curve_loop([
            model.occ.add_line(p1, p2),
            model.occ.add_line(p1, p3),
            model.occ.add_line(p2, p3)
        ])
    ])
    yoke = model.occ.cut([(2, yoke)], [(2, cut)])[0][0][1]

    nbsn_coil_1 = model.occ.add_rectangle(-487.7e-3,
                                          0, 0, nbsn_coil_w, nbsn_coil_h)
    nbsn_coil_2 = model.occ.add_rectangle(-749.7e-3,
                                          0, 0, nbsn_coil_w, nbsn_coil_h)
    nbti_coil_1 = model.occ.add_rectangle(-487.7e-3,
                                          (58 - 30) * 1e-3, 0, nbti_coil_w, nbti_coil_h)
    nbti_coil_2 = model.occ.add_rectangle(-749.7e-3,
                                          (58 - 30) * 1e-3, 0, nbti_coil_w, nbti_coil_h)
    comp_coil = model.occ.add_rectangle(-164.7e-3,
                                        0, 0, comp_coil_w, comp_coil_h)

    model.occ.translate([
        (2, yoke),
        (2, nbsn_coil_1),
        (2, nbsn_coil_2),
        (2, nbti_coil_1),
        (2, nbti_coil_2),
        (2, comp_coil),
    ], 0, yoke_oy, 0)

    air = model.occ.add_rectangle(-1000e-3, 0, 0, (1000 + 200) * 1e-3, 240e-3)
    cut = model.occ.copy([
        (2, yoke),
        (2, nbsn_coil_1),
        (2, nbsn_coil_2),
        (2, nbti_coil_1),
        (2, nbti_coil_2),
        (2, comp_coil),
    ])
    air = model.occ.cut([(2, air)], cut)[0][0][1]

    objs = [
        (2, air),
        (2, yoke),
        (2, nbsn_coil_1),
        (2, nbsn_coil_2),
        (2, nbti_coil_1),
        (2, nbti_coil_2),
        (2, comp_coil),
    ]

    tags, tag_map = model.occ.fragment(objs, [])
    model.occ.synchronize()

    mapping = np.zeros((len(tags) + 1), dtype=int)

    for i in range(len(tags)):
        mapping[objs[i][1]] = tags[i][1]

    model.add_physical_group(dim=2, tags=[mapping[air]], tag=1)
    model.add_physical_group(dim=2, tags=[mapping[yoke]], tag=2)
    model.add_physical_group(dim=2,
                             tags=[mapping[nbsn_coil_1], mapping[nbsn_coil_2]],
                             tag=3)
    model.add_physical_group(dim=2,
                             tags=[mapping[nbti_coil_1], mapping[nbti_coil_2]],
                             tag=4)
    model.add_physical_group(dim=2, tags=[mapping[comp_coil]], tag=5)

    yoke_bnd = []

    for dim, tag in model.get_boundary([(2, mapping[yoke])]):
        yoke_bnd.append(tag)

    model.add_physical_group(dim=1, tags=yoke_bnd, tag=6)

    coils_bnd = []

    for dim, tag in model.get_boundary([
        (2, mapping[nbsn_coil_1]),
        (2, mapping[nbsn_coil_2]),
        (2, mapping[nbti_coil_1]),
        (2, mapping[nbti_coil_2]),
        (2, mapping[comp_coil])
    ]):
        coils_bnd.append(tag)

    model.add_physical_group(dim=1, tags=coils_bnd, tag=7)

    bnd = model.get_boundary(
        [(2, mapping[air])], oriented=False, recursive=False)
    eps = 1e-5
    axis = []

    for dim, tag in bnd:
        if dim != 1:
            continue

        xmin, ymin, zmin, xmax, ymax, zmax = model.get_bounding_box(dim, tag)

        if abs(ymin) < eps and abs(ymax) < eps:
            axis.append(tag)

    model.add_physical_group(dim=1, tags=axis, tag=8)

    model.mesh.field.add("Distance", 1)
    model.mesh.field.set_numbers(1, "CurvesList", axis)
    model.mesh.field.set_number(1, "Sampling", 100)

    gmsh.model.mesh.field.add("Threshold", 2)
    gmsh.model.mesh.field.set_number(2, "InField", 1)
    gmsh.model.mesh.field.set_number(2, "SizeMin", 5e-3)
    gmsh.model.mesh.field.set_number(2, "SizeMax", 100e-3)
    gmsh.model.mesh.field.set_number(2, "DistMin", 30e-3)
    gmsh.model.mesh.field.set_number(2, "DistMax", 1000e-3)

    model.mesh.field.set_as_background_mesh(2)

    model.mesh.generate(dim=2)
    gmsh.write(f"{mesh_dir}/mesh.msh")
    gmsh.finalize()

    data = gmshio.read_from_msh(
        f"{mesh_dir}/mesh.msh", MPI.COMM_WORLD, gdim=2)
    plotter = pyvista.Plotter(off_screen=True)
    grid = pyvista.UnstructuredGrid(*plot.vtk_mesh(data.mesh))

    grid.cell_data["PhysicalGroups"] = data.cell_tags.values
    plotter.add_mesh(grid,
                     show_edges=True,
                     scalars="PhysicalGroups",
                     cmap="tab10",
                     opacity=0.9,
                     show_scalar_bar=False,
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

    plotter.add_mesh(
        pyvista.UnstructuredGrid(
            *plot.vtk_mesh(data.mesh, 1, data.facet_tags.find(8))),
        color="black",
        line_width=5
    )

    plotter.camera_position = "xy"
    plotter.reset_camera(bounds=[-1, 0.2, 0, 0.2, 0, 0])
    plotter.camera.zoom("tight")
    plotter.screenshot(f"{mesh_dir}/mesh.png",
                       transparent_background=True)

    plotter = pyvista.Plotter(off_screen=True)

    grid = pyvista.UnstructuredGrid(*plot.vtk_mesh(
        data.mesh, 2,
        np.concatenate((
            data.cell_tags.find(2),
            data.cell_tags.find(3),
            data.cell_tags.find(4),
            data.cell_tags.find(5)
        ))
    ))

    poly = grid.extract_surface(algorithm="dataset_surface")
    rot = poly.rotate_y(90).extrude_rotate(
        angle=270,
        resolution=30,
        capping=True
    ).rotate_y(90).rotate_x(90)
    rot = rot.cell_data_to_point_data()

    plotter.add_mesh(
        rot,
        show_edges=False,
        show_scalar_bar=False,
        cmap="plasma"
    )

    bnd = pyvista.UnstructuredGrid(*plot.vtk_mesh(
        data.mesh, 1,
        np.concatenate((
            data.facet_tags.find(6),
            data.facet_tags.find(7)
        ))
    ))

    bnd_poly = bnd.extract_surface(algorithm="dataset_surface")
    bnd_rot = bnd_poly.rotate_y(90).extrude_rotate(
        angle=270,
        resolution=30,
        capping=False
    ).rotate_y(90).rotate_x(90)
    bnd_edges = bnd_rot.extract_feature_edges(
        boundary_edges=True,
        feature_edges=True,
        non_manifold_edges=False
    )
    plotter.add_mesh(
        bnd_edges,
        color="black",
        line_width=2,
        render_lines_as_tubes=True
    )

    plotter.camera.zoom(1.3)
    plotter.screenshot(f"{mesh_dir}/3d.png", transparent_background=True)
