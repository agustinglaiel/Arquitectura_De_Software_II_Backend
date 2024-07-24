import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { loginUser, getUsers, getReservations, getHotels ,isAdmin} from './api';
import Cookies from 'js-cookie';
import '../styles/Admin.css';

const Admin = () => {
  const [reservas, setReservas] = useState([]);
  const [users, setUsers] = useState([]);
  const [hotels, setHotels] = useState([]);
  const navigate = useNavigate();
  const [errorMessage, setErrorMessage] = useState('');
  const [showError, setShowError] = useState(false);
  const [showLogin, setShowLogin] = useState(true);
  const [showOptions, setShowOption] = useState(false);
  const [showUsers, setShowUsers] = useState(false);
  const [showReservations, setShowReservations] = useState(false);
  const [showHotels, setShowHotels] = useState(false);
  
  const [formData, setFormData] = useState({
    email: '',
    password: '',
  });
  
  useEffect(() => {
    try {
      const status = isAdmin(Cookies.get("token"))
      const userData = Cookies.get('userData');
      const user = JSON.parse(userData);
      
      if (user.admin === 1) {
        setShowLogin(false);
        setShowOption(true);
        setShowError(false);
      } else {
        alert('No sos administrador');
        navigate('/');
      }
    } catch {}
  }, []);

  const handleChange = (event) => {
    const { name, value } = event.target;
    setFormData((prevFormData) => ({ ...prevFormData, [name]: value }));
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    try {
      const response = await loginUser(formData.email, formData.password);
      console.log(response.data)
      if (response.status === 200) {
        const user = {
          email: response.data.token.email,
          name: response.data.token.name,
          lastName: response.data.token.last_name,
          username: response.data.token.username,
          id: response.data.token.id,
          admin: 1,
          token: response.data.token.token
        };

        console.log(user)
        Cookies.set('userData', JSON.stringify(user));

        if (user.admin === 1) {
          setShowLogin(false);
          setShowOption(true);
          setShowError(false);
        } else {
          alert('No tenes permiso de administración');
          navigate('/');
        }
      } else if (response.status === 400) {
        setErrorMessage('El usuario no existe o la contraseña es incorrecta');
        setShowError(true);
      } else {
        setErrorMessage('Error al iniciar sesión');
        setShowError(true);
      }
    } catch (error) {
      setErrorMessage('Error al iniciar sesión');
      setShowError(true);
      console.error(error);
    }
  };

  const mostrarReservaciones = () => {
    setShowReservations(true);
    setShowHotels(false);
    setShowUsers(false);
  };

  const mostrarUsuarios = () => {
    setShowReservations(false);
    setShowHotels(false);
    setShowUsers(true);
  };

  const mostrarHoteles = () => {
    setShowHotels(true);
    setShowReservations(false);
    setShowUsers(false);
  };

  const getUsuarios = async (event) => {
    const response = await getUsers();
    const usersData = response.data.users;
    setUsers(usersData);
  };

  const getReservas = async (event) => {
    const response = await getReservations();
    const reservasData = response.data.reservations;
    setReservas(reservasData);
  };

  const getHoteles = async (event) => {
    const response = await getHotels();
    const reservasData = response.data.hotels;
    setHotels(reservasData);
  };

  return (
    <div className="admin-container">
      {showOptions && (
        <div id="Usuarios">
          <h1>Admin Panel</h1>
          <div className="tab-container">
            <div className={`tab ${showReservations ? 'active' : ''}`} onClick={() => {
              mostrarReservaciones();
              getReservas();
            }}>
              Reservaciones
            </div>
            <div className={`tab ${showUsers ? 'active' : ''}`} onClick={() => {
              mostrarUsuarios();
              getUsuarios();
            }}>
              Usuarios
            </div>
            <div className={`tab ${showHotels ? 'active' : ''}`} onClick={() => {
              mostrarHoteles();
              getHoteles();
            }}>
              Hoteles
            </div>
          </div>
          {showError && <p style={{ color: 'red' }}>{errorMessage}</p>}
        </div>
      )}

      {showLogin && (
        <div className="login-container">
          <h1 id="h1">Log In</h1>
          <form>
            <div className="form-group">
              <label htmlFor="email">Email</label>
              <input
                type="email"
                className="form-control"
                id="email"
                name="email"
                value={formData.email}
                onChange={handleChange}
              />
            </div>
            <div className="form-group">
              <label htmlFor="password">Contraseña</label>
              <input
                type="password"
                className="form-control"
                id="password"
                name="password"
                value={formData.password}
                onChange={handleChange}
              />
            </div>
            {showError && <p style={{ color: 'red' }}>{errorMessage}</p>}
            <button type="button" className="btn btn-primary" onClick={handleSubmit}>
              Iniciar sesión
            </button>
          </form>
        </div>
      )}

      {showReservations && (
        <div className="table-container">
          <table>
            <thead>
              <tr>
                <th>Id</th>
                <th>Hotel</th>
                <th>Fecha</th>
                <th>Usuario</th>
              </tr>
            </thead>
            <tbody>
              {reservas.map((reserva) => (
                <tr key={reserva.reservation_id}>
                  <td>{reserva.reservation_id}</td>
                  <td>{reserva.hotel_name}</td>
                  <td>
                    {reserva.initial_date} - {reserva.final_date}
                  </td>
                  <td>{reserva.user_name}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {showUsers && (
        <div className="table-container">
          <table>
            <thead>
              <tr>
                <th>ID</th>
                <th>Name</th>
                <th>Email</th>
              </tr>
            </thead>
            <tbody>
              {users.map((user) => (
                <tr key={user.id}>
                  <td>{user.id}</td>
                  <td>
                    {user.name} {user.lastName}
                  </td>
                  <td>{user.email}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {showHotels && (
        <div className="table-container">
          <table>
            <thead>
              <tr>
                <th>ID</th>
                <th>Nombre</th>
                <th>Habitaciones</th>
              </tr>
            </thead>
            <tbody>
              {hotels.map((hotel) => (
                <tr key={hotel.id}>
                  <td>{hotel.id}</td>
                  <td>
                    {hotel.name}
                  </td>
                  <td>{hotel.rooms_available}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
};

export default Admin;
