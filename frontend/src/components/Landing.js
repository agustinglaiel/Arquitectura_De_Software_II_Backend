import React, { useState, useEffect } from 'react';
import {Carousel} from 'react-bootstrap';
import {Link} from 'react-router-dom';
import { getImagesByHotelId, getHotels } from './api';

function Landing() {
  const [hoteles, setHoteles] = useState([]);
  const [imagenes, setImagenes] = useState([]);

  const getHoteles = async () => {
    try {
      const response = await getHotels();
      
      const hotelsData = response.data;
      console.log(response.data)
      const promises = [];
  
      for (const hotel of hotelsData) {
        promises.push(
          getImagesByHotelId(hotel.id)
            .then((response2) => {
            
              const imagesData = response2.data.images[0].Data;
              hotel.image=new Uint8Array(atob(imagesData).split('').map(char => char.charCodeAt(0)));
              return new Uint8Array(atob(imagesData).split('').map(char => char.charCodeAt(0)));
            })
            .catch((error) => {
              console.error('Error al obtener imágenes por hotel:', error);
              return null; // O maneja el error de alguna manera
            })
        );
      }
  
      const imagesArray = await Promise.all(promises);
  
      console.log(imagesArray);
      setImagenes(imagesArray.filter(image => image !== null));
      setHoteles(hotelsData);
    } catch (error) {
      console.error('Error al obtener hoteles:', error);
    }
  };

  useEffect(() => {
    getHoteles();
  }, []);
  
  return (        
  <Carousel>
    {hoteles.map((hotel) => (
      <Carousel.Item key={hotel.id} value={hotel.id}>
        {hotel.image ? (
          <img
            className="d-block w-100"
            src={`data:image/jpeg;base64,${btoa(String.fromCharCode.apply(null, hotel.image))}`}
            alt={`Imagen de ${hotel.name}`}
            style={{ objectFit: 'cover', height: '100%', maxHeight: '500px' }}
          />
        ) : (
          <img
            className="d-block w-100"
            src="https://images-wixmp-ed30a86b8c4ca887773594c2.wixmp.com/f/23e0da0f-37e6-40e5-a0d9-da5ad7056048/d5l32dr-d6117e9f-ab65-421c-8172-202a34c72455.jpg/v1/fill/w_1024,h_498,q_75,strp/goku_and_vegeta_the_kiss_by_tracexvalintyne_d5l32dr-fullview.jpg?token=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJ1cm46YXBwOjdlMGQxODg5ODIyNjQzNzNhNWYwZDQxNWVhMGQyNmUwIiwiaXNzIjoidXJuOmFwcDo3ZTBkMTg4OTgyMjY0MzczYTVmMGQ0MTVlYTBkMjZlMCIsIm9iaiI6W1t7InBhdGgiOiJcL2ZcLzIzZTBkYTBmLTM3ZTYtNDBlNS1hMGQ5LWRhNWFkNzA1NjA0OFwvZDVsMzJkci1kNjExN2U5Zi1hYjY1LTQyMWMtODE3Mi0yMDJhMzRjNzI0NTUuanBnIiwiaGVpZ2h0IjoiPD00OTgiLCJ3aWR0aCI6Ijw9MTAyNCJ9XV0sImF1ZCI6WyJ1cm46c2VydmljZTppbWFnZS53YXRlcm1hcmsiXSwid21rIjp7InBhdGgiOiJcL3dtXC8yM2UwZGEwZi0zN2U2LTQwZTUtYTBkOS1kYTVhZDcwNTYwNDhcL3RyYWNleHZhbGludHluZS00LnBuZyIsIm9wYWNpdHkiOjk1LCJwcm9wb3J0aW9ucyI6MC40NSwiZ3Jhdml0eSI6ImNlbnRlciJ9fQ.a204nMjn7GFD-WJyIyUZpyHnGXcb67GOpktYiozMUcs"
            alt="Respaldo"
            style={{ objectFit: 'cover', height: '100%', maxHeight: '500px' }}
          />  
        )}
        <Link to={`/hotel/${hotel.id}`}>
          <Carousel.Caption
            style={{
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              justifyContent: 'flex-end',
              height: '100%',
            }}
          >
            <h1>{`Maldron ${hotel.name}`}</h1>
            <p>{hotel.description}</p>
          </Carousel.Caption>
        </Link>
      </Carousel.Item>
    ))}
  </Carousel>

);
  
}

export default Landing;
