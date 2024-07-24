import React, { useState, useEffect } from 'react';
import { Container, Row, Col, Button ,Carousel, Card} from 'react-bootstrap';
import { useParams } from 'react-router-dom';
import { getHotelById, getImagesByHotelId, agregarReservation } from './api';
import { useNavigate } from 'react-router-dom';

function HotelDetail() {
  const { id, startDate, finalDate, idHabitacion } = useParams();
  const [hotel, setHotel] = useState();
  const [imagenes, setImagenes] = useState([]);
  const [showReservations, setShowReservations] = useState(false);
  const [nombreHotel,setNombre] = useState("Cargando")
  const [tipoHabitacion,setHabitacion] = useState("Cargando")
  const navigate = useNavigate();

  const getHotel = async () => {
    try {
      const response = await getHotelById(id);
      setHotel(response.data);
      console.log(response.data)
    } catch (error) {
      console.error(error);
    }
  };

  const getImagenes = async () => {
    try {
      const response2 = await getImagesByHotelId(id);
      const nose = response2.data.images
      console.log(nose)
      const imgData = nose.map(image =>
        new Uint8Array(atob(image.Data).split('').map(char => char.charCodeAt(0)))
      );
      setImagenes(imgData);
    } catch (error) {
      console.error(error);
    }
  };

  useEffect(() => {
    getHotel();
    getImagenes();
  }, [id]);

  useEffect(() => {
    if (startDate !== undefined && finalDate !== undefined ) {
      try{
        for(let i =0;i<hotel.habitaciones.length;i++){
          console.log(hotel.habitaciones[i].Id.toString())
          console.log(idHabitacion)
          if(hotel.habitaciones[i].Id.toString() === idHabitacion){
            setHabitacion(hotel.habitaciones[i].Nombre) 
          }
        }
        setNombre(hotel.name)

      }catch(error){
        console.log(error)
      }
      setShowReservations(true);
    } else {
      setShowReservations(false);
    }
  }, [startDate, finalDate, idHabitacion,hotel]);


  const handleSubmit = async (event) => {
    event.preventDefault();
    try{
      const response = await agregarReservation(
        id,
        startDate,
        finalDate
        
      );
     
      if (response.status === 200 || response.status === 201) {
        navigate(`/micuenta`);
      }else if(response.status === 401){
        navigate(`/login`);
      }else if (response.status===400) {
        console.log(`${response.data.message}`);
        alert("Fecha ya no disponible, elegi otra para continuar")
        navigate(`/reserva`);
      }
    } catch (error) {
      console.log(error)
    }
  };

  return (
    <Container className="mt-5">
      <Row>
        <Col md={6}>
          {hotel ? (
            <>
              <h1>Maldron {hotel.name}</h1>
              <p>{hotel.description}</p>
            </>
          ) : (
            <>
              <h1>Cargando...</h1>
              <p>Cargando...</p>
            </>
          )}
        </Col>
        <Col md={6}>

          
          <Carousel style={{ height: '100%', width: '100%' }}>
            {imagenes.map((imagen, index) => (
              <Carousel.Item key={index} style={{ height: '100%' }}>
                <img
                  src={`data:image/jpeg;base64,${btoa(String.fromCharCode.apply(null, imagen))}`}
                  alt={`Imagen ${index}`}
                  className="d-block w-100"
                  style={{ objectFit: 'cover', height: '100%', maxHeight: '250px' }}
                />
            </Carousel.Item>
            ))}
          </Carousel>
        </Col>
      </Row>
      <Row>
      <div className="mt-3">
        {hotel ? (
          <>
            <h5>Amenidades</h5>
            <div className="d-flex flex-wrap">
              {hotel.amenities.map((amenity, index) => (
                <Card key={index} className="m-2" style={{ width: '12rem' }}>
                  <Card.Body>
                    <Card.Text>{amenity}</Card.Text>
                  </Card.Body>
                </Card>
              ))}
            </div>
            </>
          ): (
            <>
            <h5>Amenidades</h5>
            </>
          )}

          </div>   
      </Row>
      <Row className="d-flex justify-content-between mt-4">
        {showReservations && (
          <form onSubmit={handleSubmit} className="w-100">
            <Row>
              <Col md={4} className="text-left">
                <p>Fecha de inicio: {startDate}</p>
                <p>Fecha final: {finalDate}</p>
              </Col>

              <Col md={4} className="text-center align-self-center">
                <p>Habitación: {tipoHabitacion}</p>
              </Col>

              <Col md={4} className="d-flex justify-content-end align-self-center">
                <p> </p>
                <Button type="submit" className="btn btn-primary">
                  Completar Reserva
                </Button>
              </Col>
            </Row>
          </form>
        )}
      </Row>



    </Container>
  );
}

export default HotelDetail;