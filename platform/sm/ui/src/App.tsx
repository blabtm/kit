import React, { useState } from "react";
import "./App.css";
import schema from "./model/schema.json"

import { JsonForms } from "@jsonforms/react";
import {
  materialRenderers,
  materialCells,
} from "@jsonforms/material-renderers";

import Button from "@mui/material/Button"

function App() {
  const smBaseUrl = `http://${process.env.REACT_APP_SM_URL}/v1`
  const [data, setData] = useState({})

  return (
    <div>
      <Button
        variant="contained"
        onClick={async () => {
          fetch(`${smBaseUrl}/svc/em-es/up`, {
            method: "PUT"
          })
            .then(response => {
              alert(response.text())
            })
            .catch(error => {
              console.log(error)
            })
        }}>Up</Button>

      <Button
        variant="contained"
        onClick={async () => {
          fetch(`${smBaseUrl}/svc/em-es/down`)
            .then(response => {
              alert(response.text())
            })
            .catch(error => {
              console.log(error)
            })
        }}>Down</Button>

      <Button
        variant="contained"
        onClick={async () => {
          fetch(`${smBaseUrl}/svc/em-es/config`)
            .then(response => response.json())
            .then(data => setData(data))
            .catch(error => {
              console.log(error)
            })
        }}>Pull</Button>

      <Button
        variant="contained"
        onClick={async () => {
          fetch(`${smBaseUrl}/svc/em-es/config`, {
            method: "PUT",
            headers: {
              "Content-Type": "application/json"
            },
            body: JSON.stringify(data)
          })
            .then(response => {
              alert(response.text())
            })
            .catch(error => {
              console.log(error)
            })
        }}>Push</Button>

      <JsonForms
        schema={schema}
        data={data}
        renderers={materialRenderers}
        cells={materialCells}
        onChange={({ data }) => setData(data)}
      />
    </div>
  );
}

export default App;
