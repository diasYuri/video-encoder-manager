# video-encoder-manager
video encoder manager written in golang and using library bento4 to convert mp4 into mpeg-dash



# Encoder Golang

- Recebe uma mensagem via RabbitMQ informando qual o vídeo deve ser convertido
- Faz download do vídeo no Google Cloud Storage
- Fragmenta o Vídeo
- Converte o Vídeo para MPEG-DASH
- Faz upload do vídeo no GCS
- Envia uma notificação via fila com as informações dos vídeos convertidos ou informando algum erro na conversão
- Em caso de erro, a mensagem original enviada via RabbitMQ será rejeitada e encaminhada para uma Dead Letter Exchange


