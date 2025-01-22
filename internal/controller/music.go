package controller

import (
	"log/slog"
	"music/internal/models"
	"music/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Music interface {
	GetMusics(ctx *gin.Context)
	GetSongLyricsByVerses(ctx *gin.Context)
	UpdateMusic(ctx *gin.Context)
	DeleteMusic(ctx *gin.Context)
	AddMusic(ctx *gin.Context)
}

type musicController struct {
	service service.Music
	logger  *slog.Logger
}

func newMusicController(service service.Music, logger *slog.Logger) *musicController {
	return &musicController{service: service, logger: logger}
}

// @Summary	Get musics
// @Tags		music
// @Accept		json
// @Produce	json
// @Param		group			query		string	false	"Group"
// @Param		song			query		string	false	"Song name"
// @Param		release_date	query		string	false	"Release date"
// @Param		text			query		string	false	"Text"
// @Param		offset			query		int		false	"offset"
// @Param		limit			query		int		false	"limit"
// @Success	200				{string}	Success
// @Failure	400				{string}	Bad			request
// @Failure	500				{string}	Internal	Server	Error
// @Router		/ [get]
func (c *musicController) GetMusics(ctx *gin.Context) {
	group := ctx.Query("group")
	song := ctx.Query("song")
	releaseDate := ctx.Query("release_date")
	text := ctx.Query("text")

	offset, err := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		c.logger.DebugContext(ctx, "Invalid offset", slog.String("error", err.Error()))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid offset"})
		return
	}

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		c.logger.DebugContext(ctx, "Invalid limit", slog.String("error", err.Error()))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}

	musics, err := c.service.GetMusics(ctx, group, song, releaseDate, text, limit, offset)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to get musics", slog.String("error", err.Error()))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, musics)
}

// @Summary	Get lyrics
// @Tags		music
// @Accept		json
// @Produce	json
// @Param		music_id	path		int	true	"music ID int"
// @Param		couplet		query		int	false	"couplet"
// @Param		size		query		int	false	"size"
// @Success	200			{string}	Success
// @Failure	400			{string}	Bad			request
// @Failure	500			{string}	Internal	Server	Error
// @Router		/{music_id} [get]
func (c *musicController) GetSongLyricsByVerses(ctx *gin.Context) {
	couplet, err := strconv.Atoi(ctx.DefaultQuery("couplet", "1"))
	if err != nil {
		c.logger.DebugContext(ctx, "Invalid query param page", slog.String("error", err.Error()))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid query param page"})
		return
	}

	size, err := strconv.Atoi(ctx.DefaultQuery("size", "1"))
	if err != nil {
		c.logger.DebugContext(ctx, "Invalid query param size", slog.String("error", err.Error()))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid query param size"})
		return
	}
	musicID, err := strconv.Atoi(ctx.Param("music_id"))
	if err != nil {
		c.logger.DebugContext(ctx, "Invalid ID param", slog.String("error", err.Error()))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID param"})
		return
	}

	musicText, err := c.service.GetSongLyricsByVerses(ctx, musicID, couplet, size)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to get music text", slog.String("error", err.Error()))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, musicText)
}

// @Summary	Updte musics
// @Tags		music
// @Accept		json
// @Produce	json
// @Param		music_id	path		int					true	"music ID"
// @Param		request		body		models.MusicUpdate	true	"body json"
// @Success	202			{string}	Success
// @Failure	400			{string}	Bad			request
// @Failure	500			{string}	Internal	Server	Error
// @Router		/{music_id} [patch]
func (c *musicController) UpdateMusic(ctx *gin.Context) {
	ID, err := strconv.Atoi(ctx.Param("music_id"))
	if err != nil {
		c.logger.DebugContext(ctx, "Failed to conver str to int", slog.String("error", err.Error()))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Inalid ID param"})
		return
	}

	var updates models.MusicUpdate

	if err := ctx.ShouldBindJSON(&updates); err != nil {
		c.logger.DebugContext(ctx, "Error on parsing body", slog.String("error", err.Error()))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.UpdateMusic(ctx, ID, updates); err != nil {
		c.logger.ErrorContext(ctx, "Failed to update music", slog.String("error", err.Error()))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusAccepted)
}

// @Summary	Delete music
// @Tags		music
// @Accept		json
// @Produce	json
// @Param		music_id	path		int	true	"music ID int"
// @Success	204			{string}	Success
// @Failure	400			{string}	Bad			request
// @Failure	500			{string}	Internal	Server	Error
// @Router		/{music_id} [delete]
func (c *musicController) DeleteMusic(ctx *gin.Context) {
	musicIDstr := ctx.Param("music_id")
	musicID, err := strconv.Atoi(musicIDstr)
	if err != nil {
		c.logger.DebugContext(ctx, "Error on get params", slog.String("error", err.Error()))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid music ID"})
		return
	}

	if err := c.service.DeleteMusic(ctx, musicID); err != nil {
		c.logger.ErrorContext(ctx, "Failed to delete music", slog.String("error", err.Error()))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// @Summary	Add music
// @Tags		music
// @Accept		json
// @Produce	json
// @Param		request	body		models.Music	true	"body json"
// @Success	201		{string}	Success
// @Failure	400		{string}	Bad			request
// @Failure	500		{string}	Internal	Server	Error
// @Router		/ [post]
func (c *musicController) AddMusic(ctx *gin.Context) {
	var music models.Music

	if err := ctx.ShouldBindJSON(&music); err != nil {
		c.logger.DebugContext(ctx, "Error on parse body params", slog.String("error", err.Error()))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	if err := c.service.AddMusic(ctx, music); err != nil {
		c.logger.ErrorContext(ctx, "Failed to add music", slog.String("error", err.Error()))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusCreated)
}
