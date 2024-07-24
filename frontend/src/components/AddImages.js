import React, { useState } from 'react';
import { Container, Row, Col, Button } from 'react-bootstrap';
import { useDropzone } from 'react-dropzone';
import { postImage } from './api';

function AddImages() {
  const [imagenes, setImagenes] = useState([]);
  const [previewImages, setPreviewImages] = useState([]);


  const handleEliminarImagen = (index) => {
    // Eliminar la imagen del estado
    const updatedImages = [...imagenes];
    updatedImages.splice(index, 1);
    setImagenes(updatedImages);

    // Eliminar la imagen de la previsualización
    const updatedPreview = [...previewImages];
    updatedPreview.splice(index, 1);
    setPreviewImages(updatedPreview);
  };

  const { getRootProps, getInputProps } = useDropzone({
    accept: 'image/*', // Solo permitir archivos de imagen
    onDrop: (acceptedFiles) => {
      // Actualizar el estado de las imágenes y generar la previsualización
      setImagenes([...imagenes, ...acceptedFiles]);
      const preview = [...previewImages, ...acceptedFiles.map(file => URL.createObjectURL(file))];
      setPreviewImages(preview);
    },
  });

  const handleSubmit = async (event) => {
    event.preventDefault();
    try {
      // Iterar sobre las imágenes y llamar a postImage para cada una
      for (const imagen of imagenes) {
        const response = await postImage(imagen, 2); 
        console.log(response);
      }
    } catch (error) {

      console.error('Error al procesar las imágenes:', error);
    }
  };
  

  return (
    <Container className="mt-5">
      <h2>Crear Hotel</h2>     
      <Row>
          <p>Selecciona las imagenes</p>
          <div {...getRootProps()} style={dropzoneStyles}>
            <input {...getInputProps()} />
            <p>Arrastra y suelta imágenes aquí, o haz clic para seleccionar</p>
          </div>
          <div style={previewContainerStyles}>
            {previewImages.map((url, index) => (
              <div key={index} style={{ position: 'relative' }}>
                <img
                  src={url}
                  alt={`Imagen ${index}`}
                  style={previewImageStyles}
                />
                <button
                  onClick={() => handleEliminarImagen(index)}
                  style={deleteButtonStyles}
                >
                  ×
                </button>
              </div>
            ))}
          </div>
      </Row>
      <Row className="mt-3">
        <Col md={6} className="d-flex justify-content-end">
          <Button variant="primary" onClick={handleSubmit}>
            Crear
          </Button>
        </Col>
      </Row>
    </Container>
  );
}

// Estilos para la drop zone
const dropzoneStyles = {
  border: '2px dashed #cccccc',
  borderRadius: '4px',
  padding: '20px',
  textAlign: 'center',
  cursor: 'pointer',
  marginTop: '20px',
};

// Estilos para la previsualización de imágenes
const previewContainerStyles = {
  marginTop: '20px',
  display: 'flex',
  flexWrap: 'wrap',
};

const previewImageStyles = {
  width: '100px',
  height: '100px',
  objectFit: 'cover',
  marginRight: '10px',
  marginBottom: '10px',
};

const deleteButtonStyles = {
  position: 'absolute',
  top: '5px',
  right: '5px',
  padding: '5px',
  background: 'red',
  color: 'white',
  border: 'none',
  borderRadius: '4px',
  cursor: 'pointer',
};

export default AddImages;
