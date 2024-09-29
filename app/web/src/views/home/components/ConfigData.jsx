import { useFetch } from "../../../hooks/useFetch";
import { EndpointCard } from "./EndpointCard";

export const ConfigData = () => {
  const {data, loading, error, handleCancelRequest} = useFetch("http://localhost:8080/admin/config");
  
  return (
    <div>
      {error && <div className="alert alert-danger"> {error} </div> }
      <h1 className="text-center">Current config</h1>
      {loading && 
      <div>
      <h3> cargando... </h3> 
      <button onClick={handleCancelRequest}>Cancelar request</button>
      </div>
      }
      <hr className="mt-2 mb-2" />
      {data && 
        <div className='container p-2'>
          <p>Created at: {data[0].created_at}</p>
          <hr />
          <h4>Endpoints</h4>
          <div className="d-flex d-flex-wrap">
          {
            data[0].endpoints.map(cfgEnpoint => (
              <EndpointCard endpoint={cfgEnpoint} key={cfgEnpoint}/>
            ))
          }
          </div>
        </div>
      }
    </div>
  )
}
