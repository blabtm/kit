import { useParams } from "react-router"
import { useState, useEffect, useRef } from 'react'
import { JsonForms } from '@jsonforms/react'

import AppBar from '@mui/material/AppBar'
import Button from '@mui/material/Button'
import Stack from '@mui/material/Stack'

import {
  materialRenderers,
  materialCells,
} from '@jsonforms/material-renderers'

import { dashboardRenderer, dashboardTester } from './render/Dashboard'

const customRenderers = [
  ...materialRenderers,
  { tester: dashboardTester, renderer: dashboardRenderer }
]

export function Service() {
  let params = useParams()

  const [service, setService] = useState(null)
  const [loading, setLoading] = useState(true)
  const [data, setData] = useState({})
  const [status, setStatus] = useState({})

  const intervalRef = useRef(null);

  useEffect(() => {
    async function fetchService() {
      await fetch(`http://localhost:8080/api/${params.name}/ui`, {
           mode: 'cors'
      })
        .then(r => r.json())
        .then(r => {
          let layout = JSON.parse(r.layout)

          if (Object.keys(layout).length == 0) {
            layout = undefined
          }

          let service = {
            schema: JSON.parse(r.schema),
            layout: layout
          }

          delete service.schema["$schema"]

          setService(service)
          setLoading(false)
        })
    }

    async function fetchData() {
      await fetch(`http://localhost:8080/api/${params.name}`, {
           mode: "cors"
      })
        .then(r => r.json())
        .then(r => {
          setData(r)
        })
    }

    async function fetchStatus() {
      await fetch(`http://localhost:8080/api/${params.name}/ps`, {
           mode: "cors"
      })
        .then(r => r.json())
        .then(r => {
          setStatus(r)
        })
    }

    fetchStatus()
    fetchData()
    fetchService()

    intervalRef.current = setInterval(() => {
      fetchStatus()
    }, 3000);

    return () => clearInterval(intervalRef.current);
  }, [])

  if (loading) return <p>Loading...</p>

  return (
    <div>
      <AppBar position="sticky" sx={{
        mb: "1rem"
      }}>
        <Stack
          sx={{ p: 1 }}
          alignItems="center"
          direction="row"
          spacing={2}>
          <Button
            variant="contained"
            color="secondary"
            onClick={async () => {
              await fetch(`http://localhost:8080/api/${params.name}`, {
                mode: "cors",
                method: "PUT",
                body: JSON.stringify(data),
            })
          }}>Set</Button>
          <p>{status.state}</p>
        </Stack>
      </AppBar>

      <JsonForms
        schema={service.schema}
        uischema={service.layout}
        data={data}
        renderers={customRenderers}
        cells={materialCells}
        onChange={({ data, errors }) => setData(data)}
      />
    </div>
  )
}
