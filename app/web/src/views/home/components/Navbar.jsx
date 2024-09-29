export const Navbar = () => {
  return (
    <nav className="navbar bg-dark border-bottom border-body mt-2 mb-3" data-bs-theme="dark" style={{borderRadius:"10px"}}>
  <div className="container-fluid">
    <a className="navbar-brand" >Reverse Proxy Admin</a>
    <button className="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNavAltMarkup" aria-controls="navbarNavAltMarkup" aria-expanded="false" aria-label="Toggle navigation">
      <span className="navbar-toggler-icon"></span>
    </button>
    <div className="collapse navbar-collapse" id="navbarNavAltMarkup">
      <div className="navbar-nav">
        <a className="nav-link active" aria-current="page" href="#">About</a>
        <a className="nav-link" href="#">Contact</a>
      </div>
    </div>
  </div>
</nav>
  )
}
