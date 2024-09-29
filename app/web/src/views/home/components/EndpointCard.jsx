
export const EndpointCard = ({endpoint}) => {
  return (
    <div className="card m-2" style={{width: '22rem'}} key={endpoint.prefix}>
        <div className="card-body">
            <h5 className="card-title text-center">Prefix: {endpoint.prefix}</h5>
            <h6 className="card-subtitle">Identifier: {endpoint.header_identifier}</h6>
            <ul className="list-group list-group-flush">
            {
            Object.entries(endpoint.backend_urls).map(([key, url]) => (
                <li className="list-group-item" key={key}><strong>{key}:</strong> {url}</li>
            ))
            }
            </ul>
            <a href="#" className="btn btn-primary me-1">Edit</a>
            <a href="#" className="btn btn-danger">Delete</a>
        </div>
    </div>
  )
}
