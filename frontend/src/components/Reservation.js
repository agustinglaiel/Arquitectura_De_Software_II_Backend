import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { getHotels,disponibilidadDeReserva ,getCities} from './api';

function Reservation() {

  const [cities, setCities] = useState([]);
  const [formData, setFormData] = useState({
    option1: '',
    startDate: '',
    endDate: '',
  });
  const [errorMessage, setErrorMessage] = useState('');
  const [showError, setShowError] = useState(false);

  const navigate = useNavigate();

  
  const handleChange = (event) => {
    const { name, value } = event.target;
    setFormData((prevFormData) => ({ ...prevFormData, [name]: value }));
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    try {
      if(formData.option1===""){
        setErrorMessage('Selecciona una Ciudad');
        setShowError(true);
        setCities([])
        GetCities()
        return
      }
      if(formData.startDate===""){
        setErrorMessage('Selecciona fecha inicial');
        setShowError(true);
        setCities([])
        GetCities();

        return
      }
      if(formData.endDate===""){
        setErrorMessage('Selecciona fecha final');
        setShowError(true);
        setCities([])
        GetCities();
        return
      }
      navigate(`/list/${formData.option1}/${formData.startDate}/${formData.endDate}`);
      /*
      const response = await disponibilidadDeReserva(
        formData.option1,
        formData.startDate,
        formData.endDate,
      );
      
      if (response.status === 200 || response.status === 201) {
        
      }else if(response.status === 501){
        navigate(`/registro`);

      }else if (response.status===400) {
        setErrorMessage(`${response.data.message}`);
        setShowError(true);
        setCities([])
        GetCities();
      }else{
        setCities([])
        GetCities();
        setErrorMessage(`Algo salio mal`);
        setShowError(true);
        
      }
      */
    } catch (error) {
      console.error('Error al realizar la reserva:', error);
      setCities([])
      GetCities();
      setErrorMessage('Algo salió mal');
      setShowError(true);
    }
  };

  const GetCities = async () => {
    try {
      const response = await getCities();
      setCities(response.data)
    } catch (error) {
      console.error('Error al obtener hoteles:', error);
    }
  };

  useEffect(() => {
    GetCities();
  }, []);



  return (
    <div className="container mt-5">
      <h1>Reserva</h1>
      <form onSubmit={handleSubmit}>
        <div className="row">
          <div className="col-md-12">
            <div className="form-group">
              <label htmlFor="option1">Locación</label>
              <select
                onChange={handleChange}
                value={formData.option1}
                className="form-control"
                id="option1"
                name="option1"
              >
                <option value="0">Seleccionar el lugar de su estadía</option>
                {cities.map((cities) => (
                  <option key={cities} value={cities}>
                    {cities}
                  </option>
                ))}
              </select>
            </div>
          </div>
        </div>
        <div className="row">
          <div className="col-md-6">
            <div className="form-group">
              <label htmlFor="startDate">Fecha de inicio</label>
              <input
                onChange={handleChange}
                type="date"
                className="form-control"
                id="startDate"
                name="startDate"
                value={formData.startDate}
              />
            </div>
          </div>
          <div className="col-md-6">
            <div className="form-group">
              <label htmlFor="endDate">Fecha de fin</label>
              <input
                onChange={handleChange}
                type="date"
                className="form-control"
                id="endDate"
                name="endDate"
                value={formData.endDate}
              />
            </div>
          </div>
        </div>

        {showError && <p style={{ color: 'red' }}>{errorMessage}</p>}
        <button type="submit" className="btn btn-primary">
          Buscar
        </button>
      </form>
    </div>
  );
}

export default Reservation;