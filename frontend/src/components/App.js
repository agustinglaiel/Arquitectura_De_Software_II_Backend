import React, { useState, useEffect } from 'react';
import '../styles/App.css';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import Navbar from 'react-bootstrap/Navbar';
import Nav from 'react-bootstrap/Nav';
import 'bootstrap/dist/css/bootstrap.min.css';
import Register from './Register';
import Reservation from './Reservation';
import MiCuenta from './MiCuenta';
import LogIn from './LogIn';
import Hotel from './Hotel';
import Cookies from 'js-cookie';
//import jsonServerProvider from 'ra-data-json-server';
import Admin from './Admin';
import CreateHotel from './CreateHotel';
import AddImages from './AddImages';
import Landing from './Landing';
import List from './List';

function App() {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [userData, setUserData] = useState(null);

  const userDatax = Cookies.get('userData');
  const user = userDatax ? JSON.parse(userDatax) : null;

  useEffect(() => {
    const userDataCookie = Cookies.get('userData');
    if (userDataCookie) {
      const user = JSON.parse(userDataCookie);
      console.log(user)
      setIsLoggedIn(true);
      setUserData(user);
    }
  }, []);

  const handleLogin = (data) => {
    setIsLoggedIn(true);
    setUserData(data);
  };

  const handleLogout = () => {
    setIsLoggedIn(false);
    setUserData(null);
    Cookies.remove('userData');
  };
  
  const MaldronLogo = 'https://r7c7u2r3.rocketcdn.me/wp-content/uploads/2017/12/Maldron-Hotels-and-Partners.png';

  const Footer = () => {
    return (
      <footer className="footer mt-5">
        <div className="container text-center">
          <span>© 2024 Maldron Web Page. All Rights Reserved.</span>
        </div>
      </footer>
    );
  };

  return (
    <Router>
      <div>
        <Navbar bg="light" expand="lg">
          <Navbar.Brand style={{marginLeft:'20px'}}><img src={MaldronLogo} alt="Logo"/></Navbar.Brand>
          <Navbar.Toggle aria-controls="basic-navbar-nav" />
          <Navbar.Collapse id="basic-navbar-nav">
            <Nav className="mr-auto">
              <Nav.Link as={Link} to="/">
                Inicio
              </Nav.Link>
              <Nav.Link as={Link} to="/reserva">
                Reserva
              </Nav.Link>
            </Nav>
          </Navbar.Collapse>
            <Nav className="ml-auto" style={{marginRight:'20px'}}>
              {isLoggedIn ? (
                <>
                  <Nav.Link as={Link} to="/micuenta">
                    {user.name} {user.lastName}
                  </Nav.Link>
                  <Nav.Link onClick={handleLogout}>Cerrar sesión</Nav.Link>
                </>
              ) : (
                <Nav.Link as={Link} to="/login">
                  Iniciar sesión
                </Nav.Link>
              )}
            </Nav>
        </Navbar>
        <div className="container mt-5">
          <Routes>
            <Route path="/" element={<Landing/>}/>
            <Route path="/list/:ciudad/:startDate/:finalDate" element={<List/>}/>
            <Route path="/micuenta" element={<MiCuenta usuario={userData} />}/>
            <Route path="/registro" element={<Register onLogin={handleLogin} />}/>
            <Route path="/login" element={<LogIn onLogin={handleLogin} />} />
            <Route path="/reserva" element={<Reservation />} />
            <Route path="/admin" element={<Admin />} />
            <Route path="/admin/crearHotel" element={<CreateHotel />} />
            <Route path="/admin/crearHotel/imagenes" element={<AddImages/>}/>
            <Route path="/hotel/:id" element={<Hotel/>}/>
            <Route path="/hotel/:id/:startDate/:finalDate" element={<Hotel/>}/>
          </Routes>
        </div>
        <Footer />
      </div>
    </Router>
  
  );
}

export default App;
