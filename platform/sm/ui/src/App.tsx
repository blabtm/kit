import React, { useState } from "react";
import "./App.css";
import schema from "./model/schema.json"

import { JsonForms } from "@jsonforms/react";
import {
  materialRenderers,
  materialCells,
} from "@jsonforms/material-renderers";

function App() {
  const [data, setData] = useState({})

  console.log(`service manager is at ${process.env.REACT_APP_SM_URL}`)

  return (
    <div>
      <button>Up</button>
      <button>Down</button>
      <button onClick={async () => {
        fetch(`http://${process.env.REACT_APP_SM_URL}/v1/svc/em-es/config`)
          .then(response => response.json())
          .then(data => setData(data))
          .catch(err => {
            console.log(err)
          })
      }}>Pull</button>
      <button>Push</button>

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
