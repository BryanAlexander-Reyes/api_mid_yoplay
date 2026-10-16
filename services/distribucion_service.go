
package services

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Estados operativos de la distribución (Hay que definir si la distribución será manual, automática, o si será híbrida. Es decir, que el usuario le dá al botón y se ejecuta la funcionalidad.).
type EstadoDistribucion string

const (
	Pendiente  EstadoDistribucion = "PENDIENTE"
	EnProceso  EstadoDistribucion = "EN_PROCESO"
	Generada   EstadoDistribucion = "GENERADA"
	ConError   EstadoDistribucion = "ERROR"
)

// Modalidades de distribución.
type TipoDistribucion string

const (
	PorGrupos          TipoDistribucion = "GRUPOS"
	EliminacionDirecta TipoDistribucion = "ELIMINACION_DIRECTA"
)

// Datos mínimos que necesita el servicio.
// Puedes mapear aquí tu modelo Torneo existente.
type TorneoDistribucion struct {
	ID                 int
	FechaInicio        time.Time
	Tipo               TipoDistribucion
	EstadoDistribucion EstadoDistribucion
}

// Representación mínima de un equipo inscrito.
type EquipoInscrito struct {
	ID int
}

// Respuesta que puede utilizar el controller.
type ResultadoDistribucion struct {
	IDTorneo int                `json:"id_torneo"`
	Estado   EstadoDistribucion `json:"estado_distribucion"`
	Mensaje  string             `json:"mensaje"`
}

// Contrato que debe implementar el acceso a tus CRUD.
// Los métodos de creación deben evitar registros parciales.
type DistribucionRepository interface {
	ObtenerTorneo(
		ctx context.Context,
		idTorneo int,
	) (TorneoDistribucion, error)

	ObtenerEquiposConfirmados(
		ctx context.Context,
		idTorneo int,
	) ([]EquipoInscrito, error)

	// Debe reclamar el torneo de forma atómica en la DB.
	// Devuelve false si ya está procesándose o generado.
	ReclamarDistribucion(
		ctx context.Context,
		idTorneo int,
	) (bool, error)

	MarcarEstado(
		ctx context.Context,
		idTorneo int,
		estado EstadoDistribucion,
		detalle string,
	) error

	CrearPorGrupos(
		ctx context.Context,
		torneo TorneoDistribucion,
		equipos []EquipoInscrito,
	) error

	CrearEliminacionDirecta(
		ctx context.Context,
		torneo TorneoDistribucion,
		equipos []EquipoInscrito,
	) error
}

// Servicio principal de distribución.
type DistribucionService struct {
	repo       DistribucionRepository
	minEquipos int
}

// Constructor del servicio.
func NewDistribucionService(
	repo DistribucionRepository,
) *DistribucionService {
	return &DistribucionService{
		repo:       repo,
		minEquipos: 4,
	}
}

// Generar coordina la generación completa de un torneo.
func (s *DistribucionService) Generar(
	ctx context.Context,
	idTorneo int,
) (ResultadoDistribucion, error) {
	if idTorneo <= 0 {
		return ResultadoDistribucion{},
			errors.New("el ID del torneo debe ser mayor que cero")
	}

	if s.repo == nil {
		return ResultadoDistribucion{},
			errors.New("el repositorio no está configurado")
	}

	// 1. Consultar el torneo.
	torneo, err := s.repo.ObtenerTorneo(ctx, idTorneo)
	if err != nil {
		return ResultadoDistribucion{},
			fmt.Errorf("consultando torneo: %w", err)
	}

	if torneo.ID == 0 {
		return ResultadoDistribucion{},
			fmt.Errorf("el torneo %d no existe", idTorneo)
	}

	// 2. Evitar generar una distribución dos veces.
	if torneo.EstadoDistribucion == Generada {
		return ResultadoDistribucion{
			IDTorneo: idTorneo,
			Estado:   Generada,
			Mensaje:  "La distribución ya fue generada",
		}, nil
	}

	// 3. Validar la modalidad.
	if torneo.Tipo != PorGrupos &&
		torneo.Tipo != EliminacionDirecta {
		return ResultadoDistribucion{},
			errors.New("modalidad de distribución no válida")
	}

	// 4. Reclamar la generación de forma atómica.
	// La implementación debe persistir EN_PROCESO en la DB.
	reclamado, err := s.repo.ReclamarDistribucion(
		ctx, idTorneo,
	)
	if err != nil {
		return ResultadoDistribucion{},
			fmt.Errorf("reclamando distribución: %w", err)
	}

	if !reclamado {
		return ResultadoDistribucion{
			IDTorneo: idTorneo,
			Estado:   EnProceso,
			Mensaje:  "La distribución ya está en proceso",
		}, nil
	}

	// 5. Consultar los equipos confirmados.
	equipos, err := s.repo.ObtenerEquiposConfirmados(
		ctx, idTorneo,
	)
	if err != nil {
		return s.registrarError(
			ctx, idTorneo, err,
		)
	}

	if len(equipos) < s.minEquipos {
		err := fmt.Errorf(
			"se necesitan al menos %d equipos confirmados",
			s.minEquipos,
		)
		return s.registrarError(ctx, idTorneo, err)
	}

	// 6. Generar según la modalidad del torneo.
	switch torneo.Tipo {
	case PorGrupos:
		err = s.repo.CrearPorGrupos(
			ctx, torneo, equipos,
		)

	case EliminacionDirecta:
		err = s.repo.CrearEliminacionDirecta(
			ctx, torneo, equipos,
		)
	}

	if err != nil {
		return s.registrarError(ctx, idTorneo, err)
	}

	// 7. Marcar como generada únicamente al finalizar.
	err = s.repo.MarcarEstado(
		ctx, idTorneo, Generada, "",
	)
	if err != nil {
		return ResultadoDistribucion{},
			fmt.Errorf(
				"se crearon los registros, pero no se pudo actualizar el estado: %w",
				err,
			)
	}

	return ResultadoDistribucion{
		IDTorneo: idTorneo,
		Estado:   Generada,
		Mensaje:  "Distribución generada correctamente",
	}, nil
}

// Registra el fallo para permitir su diagnóstico y recuperación.
func (s *DistribucionService) registrarError(
	ctx context.Context,
	idTorneo int,
	causa error,
) (ResultadoDistribucion, error) {
	errEstado := s.repo.MarcarEstado(
		ctx, idTorneo, ConError, causa.Error(),
	)

	if errEstado != nil {
		return ResultadoDistribucion{},
			fmt.Errorf(
				"falló la generación: %v; también falló el registro del error: %w",
				causa, errEstado,
			)
	}

	return ResultadoDistribucion{
		IDTorneo: idTorneo,
		Estado:   ConError,
		Mensaje:  "No se pudo generar la distribución",
	}, causa
}

// Calcula cuándo corresponde ejecutar la distribución:
// el día anterior al inicio, a las 00:00, hora de Colombia.

func FechaProgramadaDistribucion(
    fechaInicio time.Time,
) (time.Time, error) {
    loc, err := time.LoadLocation("America/Bogota")
    if err != nil {
        return time.Time{},
            fmt.Errorf("cargando zona horaria: %w", err)
    }

    inicio := fechaInicio.In(loc)

    medianoche := time.Date(
        inicio.Year(),
        inicio.Month(),
        inicio.Day(),
        0, 0, 0, 0,
        loc,
    )

    return medianoche.AddDate(0, 0, -1), nil
}

