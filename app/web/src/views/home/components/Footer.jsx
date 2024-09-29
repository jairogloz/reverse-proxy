export const Footer = () => {
    return (
      <div className="container fixed-bottom bg-dark" style={{borderRadius:"10px"}}>
          <footer className="py-3 my-4 bg-dark text-light" >
              <ul className="nav justify-content-center border-bottom pb-3 mb-3">
                  <li className="nav-item">
                      <a href="#" className="nav-link px-2">Home</a>
                  </li>
              </ul>
              <p className="text-center">@2024 Oscar Rodriguez</p>
          </footer>
      </div>
    )
  }