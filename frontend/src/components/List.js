import React, { useState, useEffect } from 'react';
import { Container, Row, Col, Card, Button } from 'react-bootstrap';
import { getImagesByHotelId ,querySolrCiudad,disponibilidadDeReserva} from './api';
import { useParams } from 'react-router-dom';
import { useNavigate } from 'react-router-dom';

function List() {
  const { ciudad, startDate, finalDate } = useParams();
  const [hoteles, setHoteles] = useState([]);
  const [imagenes, setImagenes] = useState([]);

  const navigate = useNavigate();

  const agregarHotel = (hotel) => {
    setHoteles((prevHoteles) => [...prevHoteles, hotel]);
  };

  
  const QuerySolr = async () => {
    const promises = [];

    try {
      if(ciudad!==undefined){
        const response = await querySolrCiudad(ciudad);
        const hotelsData = response.data;

        for (const hotel of hotelsData) {
          try{
            const response3 = await disponibilidadDeReserva(startDate,finalDate,hotel.id)
            if(response3.status()!=200){
              continue
            }
            agregarHotel(hotel) ///ACA CHATGPT
          }catch{ 
          }  
        }

      for (const hotel of hoteles) {
        promises.push(
          getImagesByHotelId(hotel.id)
           .then((response2) => {
              const imagesData = response2.data.images[0].Data;
              const imageData = new Uint8Array(atob(imagesData).split('').map(char => char.charCodeAt(0)));
              hotel.image = imageData;
              return imageData;
           })
             .catch((error) => {
             console.error('Error al obtener imágenes por hotel:', error);
              return null; // O maneja el error de alguna manera
            })
        );
      }

      const imagesArray = await Promise.all(promises);
      setImagenes(imagesArray.filter(image => image !== null));
      setHoteles(hotelsData);
    }

    } catch (error) {
      console.error('Error al obtener hoteles:', error);
    }
  };

  useEffect(() => {
    QuerySolr()

  }, []);
  
  const handleReservar = (hotelId) => {
    navigate(`/hotel/${hotelId}/${startDate}/${finalDate}`);
  };
  return (
    <Container className="mt-5">
      <h2 style={{ marginBottom: 35 }}>Resultado</h2>
      {hoteles.map((hotel) => (
        <Row key={hotel.id} className="mb-4">
          <Col md={12}>
            <Card className="h-100">
              <Row className="no-gutters">
                <Col md={4}>
                  {hotel.image && (
                    <Card.Img
                      src={`data:image/jpeg;base64,${btoa(String.fromCharCode.apply(null, hotel.image))}`}
                      alt={`Imagen de ${hotel.name}`}
                      className="img-fluid"
                      style={{ maxHeight: '200px', objectFit: 'cover' }}
                    />
                  )}
                </Col>
                <Col md={8}>
                  <Card.Body className="d-flex flex-column">
                    <Row>
                      <Col md={4}>
                        <Card.Title>{hotel.name}</Card.Title>
                      </Col>
                      <Col md={8} />
                    </Row>
                    <Row>
                      <Col md={12}>
                        <Card.Text>{hotel.description}</Card.Text>
                      </Col>
                    </Row>
                    <Row className="mt-auto">
                      <Col md={12} className="d-flex justify-content-end">
                        <Button variant="primary" onClick={() => handleReservar(hotel.id)}>Reservar</Button>
                      </Col>
                    </Row>
                  </Card.Body>
                </Col>
              </Row>
            </Card>
          </Col>
        </Row>
      ))}
    </Container>
  );
}

export default List;
