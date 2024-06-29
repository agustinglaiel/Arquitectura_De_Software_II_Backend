FROM ubuntu:20.04

# Establecer las variables de entorno para interactuar con los repositorios no interactivos
ENV DEBIAN_FRONTEND=noninteractive

# Actualizar el sistema e instalar dependencias necesarias
RUN apt-get update && apt-get install -y \
    software-properties-common \
    curl \
    gnupg \
    apt-transport-https

# Añadir el repositorio de RabbitMQ
RUN curl -fsSL https://packagecloud.io/rabbitmq/rabbitmq-server/gpgkey | apt-key add -
RUN echo "deb https://packagecloud.io/rabbitmq/rabbitmq-server/ubuntu/ focal main" | tee /etc/apt/sources.list.d/rabbitmq-server.list

# Añadir el repositorio de Erlang Solutions, necesario para RabbitMQ
RUN curl -fsSL https://packages.erlang-solutions.com/ubuntu/erlang_solutions.asc | apt-key add -
RUN echo "deb https://packages.erlang-solutions.com/ubuntu focal contrib" | tee /etc/apt/sources.list.d/erlang.list

# Actualizar e instalar RabbitMQ y sus dependencias
RUN apt-get update && apt-get install -y \
    erlang-base \
    erlang-asn1 \
    erlang-crypto \
    erlang-eldap \
    erlang-ftp \
    erlang-inets \
    erlang-mnesia \
    erlang-os-mon \
    erlang-parsetools \
    erlang-public-key \
    erlang-runtime-tools \
    erlang-snmp \
    erlang-ssl \
    erlang-syntax-tools \
    erlang-tftp \
    erlang-tools \
    erlang-xmerl \
    rabbitmq-server=3.13.3-1

# Exponer los puertos para RabbitMQ y su consola de administración
EXPOSE 5672 15672

# Comando para iniciar RabbitMQ
CMD ["rabbitmq-server"]

