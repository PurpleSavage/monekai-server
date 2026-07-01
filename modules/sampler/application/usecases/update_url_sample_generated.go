package samplerusecases

import (
	"fmt"
	samplerequestsdto "github.com/PurpleSavage/monekai-server/modules/sampler/application/dtos/requests"
	samplerports "github.com/PurpleSavage/monekai-server/modules/sampler/application/ports"
	samplervalueobjects "github.com/PurpleSavage/monekai-server/modules/sampler/domain/valueobjects"
	commonports "github.com/PurpleSavage/monekai-server/modules/shared/common/application/ports"
	commonservices "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/services"
	"github.com/google/uuid"
)

type UpdateUrlSampleGenerated struct {
	audioPersistenceService samplerports.SamplerPersistencePort
	storageService commonports.StoragePort
}

func NewUpdateUrlSampleGenerated(
	audioPersistenceService samplerports.SamplerPersistencePort,
	storageService commonports.StoragePort,

) *UpdateUrlSampleGenerated {
	return &UpdateUrlSampleGenerated{
		audioPersistenceService:audioPersistenceService,
		storageService:storageService,
	}
}

func (u *UpdateUrlSampleGenerated) Execute(
	data samplerequestsdto.SongWebhookResponse,
) (*samplervalueobjects.ResponseUpdateurlVo, error) {
	
	// 1. Descargar el audio temporal desde Replicate
	ioReader,size, err := commonservices.DownloadFile(data.Output)
	if err != nil {
		return nil, err 
	}
	defer ioReader.Close()

	// 2. Generar un nombre único y seguro para el bucket (evita el bug de caracteres raros o espacios de la UI)
	name := fmt.Sprintf("%s.mp3", uuid.NewString()) 
	path := fmt.Sprintf("samples/%s", name)

	// 3. Subir el archivo a tu servicio de Storage
	url, err := u.storageService.UploadFile(path, ioReader,size, "audio/mpeg")
	if err != nil {
		return nil, err
	}

	// 4. Actualizar en la base de datos y recuperar el VO con los datos del usuario + la entidad
	responseVo, err := u.audioPersistenceService.UpdateAudioUrl(url, data.ID)
	if err != nil {
		return nil, err
	}

	return responseVo, nil
}