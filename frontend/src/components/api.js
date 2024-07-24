import axios from "axios";
import Cookies from "js-cookie";

const URL_API_USER = "http://localhost:8060";
 const URL_API_SEAR = "http://localhost:8070";
 const URL_API_FICH = "http://localhost:8080";

//Register

export const postUser = async (
  FirstName,
  LastName,
  Username,
  Email,
  Password
) => {
  try {
    const response = await axios.post(`${URL_API_USER}/register`, {
      first_name: FirstName,
      last_name: LastName,
      username: Username,
      email: Email,
      password: Password,
    });
    return response;
  } catch (error) {
    if (error.response.status === 400) {
      // El servidor respondió con un código de estado de error
      const errorMessage = error.response.data;
      // Manejar el mensaje de error, por ejemplo, mostrarlo en la interfaz de usuario
      console.error(errorMessage);
      return error.response;
    } else {
      // Error de red o solicitud cancelada
      console.error("Error en la solicitud:", error.message);
    }
  }

  throw new Error("Error al registrar el usuario");
};

//Hotel

export const postHotel = async (name, Nroom, descr) => {
  try {
    
    const response = await axios.post(
      `${URL_API_USER}/insertHhotel/${name}/${Nroom}/${descr}`
    );
    return response;
  } catch (error) {
    if (error.response.status === 400) {
      // El servidor respondió con un código de estado de error
      const errorMessage = error.response.data;
      // Manejar el mensaje de error, por ejemplo, mostrarlo en la interfaz de usuario
      console.error(errorMessage);
      return error.response;
    } else {
      // Error de red o solicitud cancelada
      console.error("Error en la solicitud:", error.message);
    }
  }

  throw new Error("Error al agregar hotel");
};

//Imagen
export const postImage = async (image, idHotel) => {
  try {
    const response = await axios.post(
      `${URL_API_USER}/postImage/${image}/${idHotel}`
    );
    return response;
  } catch (error) {
    if (error.response.status === 400) {
      // El servidor respondió con un código de estado de error
      const errorMessage = error.response.data;
      // Manejar el mensaje de error, por ejemplo, mostrarlo en la interfaz de usuario
      console.error(errorMessage);
      return error.response;
    } else {
      // Error de red o solicitud cancelada
      console.error("Error en la solicitud:", error.message);
    }
  }

  throw new Error("Error al agregar imagenes");
};

export const loginUser = async (username, password) => {
  const data = {
    username: username,
    password: password,
  };

  try {
    const response = await axios.post(`${URL_API_USER}/login`, data);
    return response;
  } catch (error) {
    if (error.response.status === 400) {
      // El servidor respondió con un código de estado de error
      const errorMessage = error.response.data;
      // Manejar el mensaje de error, por ejemplo, mostrarlo en la interfaz de usuario
      console.error(errorMessage);
      return error.response;
    } else {
      // Error de red o solicitud cancelada
      console.error("Error en la solicitud:", error.message);
    }
  }
};

export const isAdmin=async(token) =>{
  try {

    const userDataFromCookie=Cookies.get("userData")
    const userData = JSON.parse(userDataFromCookie);
    axios.defaults.headers.common["Authorization"] = userData.token;
  } catch (error) {
    alert("Token expirado, registrece de nuevo para continuar");
  }
    try{
      const response = await axios.head(`${URL_API_USER}/admin/`);
    return response;}
    catch(error){
      console.log(error)
    }
  
} 


//Reservation
export const agregarReservation = async (
  idHotel,
  inicio,
  final,
  habitacion
) => {
  try {
    const userData = Cookies.get("userData");
    const user = JSON.parse(userData);
    axios.defaults.headers.common["Authorization"] = user.token;
  } catch (error) {
    alert("Token expirado, registrece de nuevo para continuar");
  }
  try {
    const response = await axios.post(
      `${URL_API_USER}/usuario/agregarReservation/${idHotel}/${inicio}/${final}/${habitacion}`
    );
    return response;
  } catch (error) {
    return error.response;
  }
};

//Dispponibilidad de reservas
export const disponibilidadDeReserva = async (
  inicio,
  final,
  idHotel
) => {
  try {
    const response = await axios.get(
      `${URL_API_USER}/reserva/${inicio}/${final}/${idHotel}`
    );
    if (response.status === 200 || response.status === 201) {
      try {
        const userData = Cookies.get("userData");
        const user = JSON.parse(userData);
        console.log(user.token);
      } catch (error) {
        alert("Reservacion Disponible, registrece para continuar");
        response.status = 501;
        return response;
      }
      return response;
    }
  } catch (error) {
    console.log(error);
    return error.response;
  }
};

//Get users
export const getUsers = async () => {
  try {
    const userData = Cookies.get("userData");
    const user = JSON.parse(userData);

    axios.defaults.headers.common["Authorization"] = user.token;
    const response = await axios.get(`${URL_API_USER}/admin/users`);
    return response;
  } catch (error) {
    console.error("Error al obtener los usuarios:", error);
  }
};

export const getCities = async () => {
  try {
    const response = await axios.get(`${URL_API_SEAR}/ciudades`);
    console.log("si me llama el get cities")
    return response;
  } catch (error) {
    console.error("Error al obtener las ciudades", error);
  }
};


export const querySolrCiudad = async(ciudad)=> {
  try {
    const response = await axios.get(`${URL_API_SEAR}/busqueda_hotel_api/search=city_${ciudad}`);
    return response
}catch(response){
  console.log("No ciudad");
}

}


//GetHotels
export const getHotels = async () => {
  try {
    const userData = Cookies.get("userData");
    const user = JSON.parse(userData);
    axios.defaults.headers.common["Authorization"] = user.token;
  } catch (error) {
    console.log("Usuario no regsitrado");
  }
  try {
    const response = await axios.get(`${URL_API_FICH}/hotels`);
    
    return response;
  } catch (error) {
    console.error("Error al obtener los usuarios:", error);
  }
};

export const getHotelById = async (id) => {
  try {
    const response = await axios.get(`${URL_API_FICH}/hotel/${id}`);
    return response;
  } catch (error) {
    console.error("Error al obtener los usuarios:", error);
  }
};

export const getReservations = async () => {
  try {
    const userData = Cookies.get("userData");
    const user = JSON.parse(userData);

    axios.defaults.headers.common["Authorization"] = user.token;

    const response = await axios.get(`${URL_API_USER}/admin/reservas`);
    return response;
  } catch (error) {
    console.error("Error al obtener los usuarios:", error);
  }
};

export const getReservationsByUser = async () => {
  try {
    const userData = Cookies.get("userData");
    const user = JSON.parse(userData);

    axios.defaults.headers.common["Authorization"] = user.token;
    const id = user.id;

    const response = await axios.get(
      `${URL_API_USER}/usuario/reservaByUserId/${id}`
    );
    console.log(response.data);
    return response;
  } catch (error) {
    console.error("Error al obtener los usuarios:", error);
  }
};

export const getImagesByHotelId = async (idHotel) => {
  try {
    const response = await axios.get(
      `${URL_API_FICH}/hotel/${idHotel}/images`
    );
    return response;
  } catch (error) {
    console.error("Error al obtener las imagenes:", error);
  }
};



export const InsertHotel = async (data) => {
  try {
    const userData = Cookies.get("userData");
    const user = JSON.parse(userData);

    axios.defaults.headers.common["Authorization"] = user.token;

    const response = await axios.post(
      `${URL_API_USER}admin/InsertHotel`,
      data,
      {
        headers: {
          "Content-Type": "application/json",
          // Otros encabezados si es necesario
        },
      }
    );

    return response.data; // Retorna los datos de la respuesta
  } catch (error) {
    console.error("Error al insertar el hotel:", error);
    throw error; // Lanza el error para que pueda ser manejado por el código que llama a InsertHotel
  }
};

export const tipoHabitaciones = async () => {
  try {
    const userData = Cookies.get("userData");
    const user = JSON.parse(userData);

    axios.defaults.headers.common["Authorization"] = user.token;
    const response = await axios.get(`${URL_API_USER}/admin/Habitaciones`);
    return response;
  } catch (error) {
    console.error("Error al obtener los usuarios:", error);
  }
};

export const getAmenities = async () => {
  try {
    const response = await axios.get(`${URL_API_USER}/getAmenities`);
    return response;
  } catch (error) {
    console.error("Error al obtener los usuarios:", error);
  }
};
