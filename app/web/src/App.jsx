import { useFetch } from './hooks/useFetch'
import './App.css'

function App() {

  const {data, loading, error, handleCancelRequest} = useFetch("http://localhost:8080/admin/config");

  return (
    <>
    <div className='container'>
      <h1>Fetch like a fucking pro</h1>
      {loading && 
      <div>
      <h3> cargando... </h3> 
      <button onClick={handleCancelRequest}>Cancelar request</button>
      </div>
      }

      {error && <h3> ocurrio un error </h3> }
      <div className='card'>
        {JSON.stringify(data)}
      </div>
    </div>
      
    </>
  )
}

export default App
