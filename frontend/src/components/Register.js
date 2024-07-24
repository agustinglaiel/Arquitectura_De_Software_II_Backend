import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { postUser } from "./api";
import Cookies from "js-cookie";

const Register = ({ onLogin }) => {
  const navigate = useNavigate();

  const [formData, setFormData] = useState({
    firstName: "",
    lastName: "",
    userName: "",
    email: "",
    password: "",
    confirmPassword: "",
  });

  const [errorMessage, setErrorMessage] = useState("");
  const [showError, setShowError] = useState(false);

  const handleChange = (event) => {
    const { name, value } = event.target;
    setFormData((prevFormData) => ({ ...prevFormData, [name]: value }));
  };

  const handleSubmit = async (event) => {
    event.preventDefault();

    if (formData.password !== formData.confirmPassword) {
      // Las contraseñas no coinciden, muestra un mensaje de error
      setErrorMessage("Las contraseñas no coinciden");
      setShowError(true);
      return;
    }

    try {
      const response = await postUser(
        formData.firstName,
        formData.lastName,
        formData.userName,
        formData.email,
        formData.password
      );

      if (response.status === 201) {
        console.log(response);
        const user = {
          email: response.data.email,
          firstName: response.data.firstName,
          lastName: response.data.lastName,
          userName: response.data.userName,
          id: response.data.id,
          token: response.data.token,
        };

        Cookies.set("userData", JSON.stringify(user));
        onLogin(formData.firstName, formData); // Llama a la función onLogin pasando el nombre del usuario registrado y los datos del formulario
        navigate("/"); // Redirige al usuario a la página principal después de registrar exitosamente
      } else if (response.status === 400) {
        setErrorMessage("El correo electrónico o DNI ya está en uso");
        setShowError(true);
        return;
      } else {
        setErrorMessage("Error al registrar el usuario");
      }
    } catch (error) {}

    // Restablece los valores y oculta el mensaje de error

    setFormData({
      firstName: "",
      lastName: "",
      userName: "",
      email: "",
      password: "",
      confirmPassword: "",
    });

    setShowError(false);
  };

  return (
    <div className="container">
      <h1 id="h1">Sign In</h1>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="firstName">Nombre</label>
          <input
            type="text"
            className="form-control"
            id="firstName"
            name="firstName"
            value={formData.firstName}
            onChange={handleChange}
          />
        </div>
        <div className="form-group">
          <label htmlFor="lastName">Apellido</label>
          <input
            type="text"
            className="form-control"
            id="lastName"
            name="lastName"
            value={formData.lastName}
            onChange={handleChange}
          />
        </div>
        <div className="form-group">
          <label htmlFor="userName">Username</label>
          <input
            type="text"
            className="form-control"
            id="userName"
            name="userName"
            value={formData.userName}
            onChange={handleChange}
          />
        </div>
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
        <div className="form-group">
          <label htmlFor="confirmPassword">Confirmar Contraseña</label>
          <input
            type="password"
            className="form-control"
            id="confirmPassword"
            name="confirmPassword"
            value={formData.confirmPassword}
            onChange={handleChange}
          />
        </div>
        {showError && <p style={{ color: "red" }}>{errorMessage}</p>}
        <button type="submit" className="btn btn-primary">
          Registrarse
        </button>
      </form>
    </div>
  );
};

export default Register;
