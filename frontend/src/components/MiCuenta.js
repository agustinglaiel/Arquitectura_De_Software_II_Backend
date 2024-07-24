import '../styles/MiCuenta.css';
import Cookies from 'js-cookie';
import { useNavigate } from 'react-router-dom';

const MiCuenta = () => {
  const userData = Cookies.get('userData');
  const navigate = useNavigate();

  if (!userData) {
    navigate('/');
    return (
      <div className="container">
        <h1>Mi Cuenta</h1>
        <p>No se encontraron datos de usuario.</p>
      </div>
    );
  }

  const user = JSON.parse(userData);

  return (
    <div className="container">
      <h1>Mi Cuenta</h1>
      <div className="user-details">
        <div className="user-image">
          
            <img
              id="perfil"
              src="https://static.vecteezy.com/system/resources/previews/002/387/693/large_2x/user-profile-icon-free-vector.jpg"
              alt="Foto de perfil"
              width="250"
              height="250"
            />
         
        </div>
        <div className="user-info">
          <p className="user-info-line">
            <span className="label">Nombre:</span> <span className="value">{user.name}</span>
          </p>
          <p className="user-info-line">
            <span className="label">Apellido:</span> <span className="value">{user.lastName}</span>
          </p>
          <p className="user-info-line">
            <span className="label">Correo electrónico:</span> <span className="value">{user.email}</span>
          </p>
          <p className="user-info-line">
            <span className="label">DNI:</span> <span className="value">{user.dni}</span>
          </p>
        </div>
        
      </div>

    </div>
  );
};

export default MiCuenta;
