package samplerusecases

import (
	"log"

	samplerequestsdto "github.com/PurpleSavage/monekai-server/modules/sampler/application/dtos/requests"
	samplerports "github.com/PurpleSavage/monekai-server/modules/sampler/application/ports"
	samplerentities "github.com/PurpleSavage/monekai-server/modules/sampler/domain/entities"
	samplervalueobjects "github.com/PurpleSavage/monekai-server/modules/sampler/domain/valueobjects"
	samplerutils "github.com/PurpleSavage/monekai-server/modules/sampler/utils"
)


type GenerateSampleUseCase struct{
	songGeneratorService samplerports.SongGeneratorPort
	audioRepository samplerports.SamplerPersistencePort
}

func NewGeneratorSampleUseCase(
	songGeneratorService samplerports.SongGeneratorPort,
	audioRepository samplerports.SamplerPersistencePort,
)*GenerateSampleUseCase{
	return  &GenerateSampleUseCase{
		songGeneratorService:songGeneratorService,
		audioRepository: audioRepository,
	}
}

func (g *GenerateSampleUseCase) Execute(
	dto samplerequestsdto.GenerateSampleDTO, 
	userId string,
)(*samplerentities.SongGeneratorResponse, error){

	vo, appErr := samplervalueobjects.CreatePayloadGenerateSongVO(
		samplerutils.BuildWebhook("songs"),
		dto.ModelVersion,
		dto.Prompt,
		dto.Duration,
		dto.OutputFormat,
		userId,
	)
	if appErr != nil{
		log.Printf("error del value obect payload: %v", appErr)
		return  nil, appErr
	}
	response, err:= g.songGeneratorService.GenerateSong(vo)
	if err != nil {
		log.Printf("error del song generator: %v", err)
		return  nil,err
	}
	saveSampleVo,voErr :=samplervalueobjects.CreateSaveSampleVO(
		userId,
		dto.SampleName,
		dto.Prompt,
		dto.Duration,
		string(dto.ModelVersion),
		string(dto.OutputFormat),
		response.GenerationId,
		string(response.StatusGeneration),
	)
	if voErr!= nil{
		log.Printf("error del value obect create-sample: %v", voErr)
		return  nil, voErr
	}
	err = g.audioRepository.SaveSampleWithoutAudioUrl(*saveSampleVo)
	if err != nil {
		log.Printf("error de la bd: %v", voErr)
		return nil, err
	}
	return  response,nil
}