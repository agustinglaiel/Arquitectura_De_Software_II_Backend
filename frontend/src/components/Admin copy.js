import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { loginUser, getUsers, getReservations } from './api';
import Cookies from 'js-cookie';
import '../styles/Admin.css';

const Admin = () => {
  const [reservas, setReservas] = useState([]);
  const [users, setUsers] = useState([]);
  const navigate = useNavigate();
  const [errorMessage, setErrorMessage] = useState('');
  const [showError, setShowError] = useState(false);
  const [showLogin, setShowLogin] = useState(true);
  const [showOptions, setShowOption] = useState(false);
  const [showUsers, setShowUsers] = useState(false);
  const [showReservations, setShowReservations] = useState(false);
  
  const [formData, setFormData] = useState({
    email: '',
    password: '',
  });

  useEffect(() => {
    try {
      const userData = Cookies.get('userData');
      const user = JSON.parse(userData);
      if (user.admin === 1) {
        setShowLogin(false);
        setShowOption(true);
        setShowError(false);
      } else {
        alert('No sos administrador');
        Cookies.remove('userData')
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
      if (response.status === 200) {
        const user = {
          email: response.data.email,
          name: response.data.name,
          lastName: response.data.lastName,
          dni: response.data.dni,
          id: response.data.id,
          admin: response.data.admin,
          token: response.data.token
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
    setShowUsers(false);
  };

  const mostrarUsuarios = () => {
    setShowReservations(false);
    setShowUsers(true);
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
    </div>
  );
};

export default Admin;
