import React, { useEffect, useState } from 'react';
import { ToastContainer, toast } from "react-toastify";

const notifyCreated = () => {
  toast.success("Contenedor creado!", {
      pauseOnHover: false,
      autoClose: 2000,
  })
}

const notifyCreationError = () => {
  toast.error("Hubo un error al crear el contenedor!", {
      pauseOnHover: false,
      autoClose: 2000,
  })
}

const notifyStopped = () => {
  toast.success("Contenedor parado!", {
      pauseOnHover: false,
      autoClose: 2000,
  })
}

const notifyStoppingError = () => {
  toast.error("Hubo un error al parar el contenedor!", {
      pauseOnHover: false,
      autoClose: 2000,
  })
}

const notifyRemove = () => {
  toast.success("Contenedor elimnado!", {
      pauseOnHover: false,
      autoClose: 2000,
  })
}

const notifyRemotionError = () => {
  toast.error("Hubo un error al eleimar el contenedor!", {
      pauseOnHover: false,
      autoClose: 2000,
  })
}

const notifyStarted = () => {
  toast.success("Contenedor iniciado!", {
      pauseOnHover: false,
      autoClose: 2000,
  })
}

const notifyStartError = () => {
  toast.error("Hubo un error al iniciar el contenedor!", {
      pauseOnHover: false,
      autoClose: 2000,
  })
}


const AdminInfra = () => {
    const [contenedores, setContenedores] = useState([]);
    
    const getContenedores = async () => {
      try {
        const request = await fetch("http://localhost:8040/containers");
        const response = await request.json();
        console.log("Contenedores obtenidos:", response);  // Para ver los datos obtenidos en la consola
        setContenedores(response);
      } catch (error) {
        console.log("No se pudieron obtener los contenedores:", error);
      }
    };
  
    useEffect(() => {
      getContenedores();
    }, []);
  
    const Cuenta = () => {
      window.location.href = '/micuenta';
    };
  
    const Home = () => {
      window.location.href = '/';
    };
  
    const handleVolver = () => {
      window.history.back();
    };
  
    const handleApagar = (id) => {
      console.log("Apagar contenedor con ID:", id);
      // Lógica para apagar el contenedor
    };
  
    const handlePrender = (id) => {
      console.log("Prender contenedor con ID:", id);
      // Lógica para prender el contenedor
    };
  
    const handleBorrar = (id) => {
      console.log("Borrar contenedor con ID:", id);
      // Lógica para borrar el contenedor
    };
  
    const handleCrear = (image, name, label, id) => {
      console.log("Crear nueva instancia de contenedor:", { image, name, label, id });
      // Lógica para crear una nueva instancia del contenedor
    };
  
    return (
      <div className="bodyinicio">
        <div className="header-content-infra">
          <div className="admin-button-container">
            <button className="admin-button" onClick={Home}>
              Inicio
            </button>
          </div>
          <div className="cuenta-button-container">
            <button className="cuenta-button" onClick={Cuenta}>
              Tu Cuenta
            </button>
          </div>
          <div className="admin-button-container">
            <button className="admin-button" onClick={handleVolver}>
              Volver
            </button>
          </div>
        </div>
        <div className="containerIni">
          <div className="hotels-containerH">
            {contenedores.length ? (
              contenedores.map((contenedor) => (
                <div className="hotel-cardH" key={contenedor.Id}>
                  <div className='img-name'>
                    <div className="hotel-infoH">
                      <h4> Contenedor: {contenedor.Names[0]} </h4>
                      <div className="hotel-descriptionH">
                        <label htmlFor={`description-${contenedor.Id}`}> Imagen: {contenedor.Image} </label>
                        <label htmlFor={`description-${contenedor.Id}`}> Estado: {contenedor.State}</label>
                      </div>
                      {contenedor.State !== "exited" && (
                        <button className="botonAC" onClick={() => handleApagar(contenedor.Id)}>Apagar</button>
                      )}
                      {contenedor.State === "exited" && (
                        <button className="botonAC" onClick={() => handlePrender(contenedor.Id)}>Prender</button>
                      )}
                      <button className="botonAC" onClick={() => handleBorrar(contenedor.Id)}>Borrar</button>
                      <button className="botonAC" onClick={() => handleCrear(contenedor.Image, contenedor.Names[0], contenedor.Labels["org.opencontainers.image.ref.name"], contenedor.Id)}> Crear nueva instancia </button>
                    </div>
                  </div> 
                </div>
              ))
            ) : (
              <p>No hay contenedores disponibles</p>
            )}
          </div>
        </div>
        <ToastContainer />
      </div>
    );
  };
  
  export default AdminInfra;