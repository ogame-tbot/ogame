package ogame

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"

	"github.com/labstack/echo"
)

// APIResp ...
type APIResp struct {
	Status  string
	Code    int
	Message string
	Result  interface{}
}

// SuccessResp ...
func SuccessResp(data interface{}) APIResp {
	return APIResp{Status: "ok", Code: 200, Result: data}
}

// ErrorResp ...
func ErrorResp(code int, message string) APIResp {
	return APIResp{Status: "error", Code: code, Message: message}
}

// HomeHandler ...
func HomeHandler(c echo.Context) error {
	version := c.Get("version").(string)
	commit := c.Get("commit").(string)
	date := c.Get("date").(string)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"version": version,
		"commit":  commit,
		"date":    date,
	})
}

// TasksHandler return how many tasks are queued in the heap.
func TasksHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	return c.JSON(http.StatusOK, SuccessResp(bot.GetTasks()))
}

// GetServerHandler ...
func GetServerHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	return c.JSON(http.StatusOK, SuccessResp(bot.GetServer()))
}

// GetServerDataHandler ...
func GetServerDataHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	return c.JSON(http.StatusOK, SuccessResp(bot.serverData))
}

// SetUserAgentHandler ...
// curl 127.0.0.1:1234/bot/set-user-agent -d 'userAgent="New user agent"'
func SetUserAgentHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	userAgent := c.Request().PostFormValue("userAgent")
	bot.SetUserAgent(userAgent)
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// ServerURLHandler ...
func ServerURLHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	return c.JSON(http.StatusOK, SuccessResp(bot.ServerURL()))
}

// GetLanguageHandler ...
func GetLanguageHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	return c.JSON(http.StatusOK, SuccessResp(bot.GetLanguage()))
}

// PageContentHandler ...
// curl 127.0.0.1:1234/bot/page-content -d 'page=overview&cp=123'
func PageContentHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	if err := c.Request().ParseForm(); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	pageHTML, _ := bot.GetPageContent(c.Request().Form)
	return c.JSON(http.StatusOK, SuccessResp(pageHTML))
}

// LoginHandler ...
func LoginHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	if _, err := bot.LoginWithExistingCookies(); err != nil {
		if err == ErrBadCredentials {
			return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// LogoutHandler ...
func LogoutHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	bot.Logout()
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// GetUsernameHandler ...
func GetUsernameHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	return c.JSON(http.StatusOK, SuccessResp(bot.GetUsername()))
}

// GetUniverseNameHandler ...
func GetUniverseNameHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	return c.JSON(http.StatusOK, SuccessResp(bot.GetUniverseName()))
}

// GetUniverseSpeedHandler ...
func GetUniverseSpeedHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	return c.JSON(http.StatusOK, SuccessResp(bot.serverData.Speed))
}

// GetUniverseSpeedFleetHandler ...
func GetUniverseSpeedFleetHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	return c.JSON(http.StatusOK, SuccessResp(bot.serverData.SpeedFleet))
}

// ServerVersionHandler ...
func ServerVersionHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	return c.JSON(http.StatusOK, SuccessResp(bot.serverData.Version))
}

// ServerTimeHandler ...
func ServerTimeHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	return c.JSON(http.StatusOK, SuccessResp(bot.ServerTime()))
}

// IsUnderAttackHandler ...
func IsUnderAttackHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	isUnderAttack, err := bot.IsUnderAttack()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(isUnderAttack))
}

// IsVacationModeHandler ...
func IsVacationModeHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	isVacationMode := bot.isVacationModeEnabled
	return c.JSON(http.StatusOK, SuccessResp(isVacationMode))
}

// GetUserInfosHandler ...
func GetUserInfosHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	return c.JSON(http.StatusOK, SuccessResp(bot.GetUserInfos()))
}

// GetCharacterClassHandler ...
func GetCharacterClassHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	return c.JSON(http.StatusOK, SuccessResp(bot.CharacterClass()))
}

// HasCommanderHandler ...
func HasCommanderHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	hasCommander := bot.hasCommander
	return c.JSON(http.StatusOK, SuccessResp(hasCommander))
}

// HasAdmiralHandler ...
func HasAdmiralHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	hasAdmiral := bot.hasAdmiral
	return c.JSON(http.StatusOK, SuccessResp(hasAdmiral))
}

// HasEngineerHandler ...
func HasEngineerHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	hasEngineer := bot.hasEngineer
	return c.JSON(http.StatusOK, SuccessResp(hasEngineer))
}

// HasGeologistHandler ...
func HasGeologistHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	hasGeologist := bot.hasGeologist
	return c.JSON(http.StatusOK, SuccessResp(hasGeologist))
}

// HasTechnocratHandler ...
func HasTechnocratHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	hasTechnocrat := bot.hasTechnocrat
	return c.JSON(http.StatusOK, SuccessResp(hasTechnocrat))
}

// GetEspionageReportMessagesHandler ...
func GetEspionageReportMessagesHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	report, err := bot.GetEspionageReportMessages()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(report))
}

// GetEspionageReportHandler ...
func GetEspionageReportHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	msgID, err := strconv.ParseInt(c.Param("msgid"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid msgid id"))
	}
	espionageReport, err := bot.GetEspionageReport(msgID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(espionageReport))
}

// GetEspionageReportForHandler ...
func GetEspionageReportForHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	galaxy, err := strconv.ParseInt(c.Param("galaxy"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid galaxy"))
	}
	system, err := strconv.ParseInt(c.Param("system"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid system"))
	}
	position, err := strconv.ParseInt(c.Param("position"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid position"))
	}
	planet, err := bot.GetEspionageReportFor(Coordinate{Type: PlanetType, Galaxy: galaxy, System: system, Position: position})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(planet))
}

// SendMessageHandler ...
// curl 127.0.0.1:1234/bot/send-message -d 'playerID=123&message="Sup boi!"'
func SendMessageHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	playerID, err := strconv.ParseInt(c.Request().PostFormValue("playerID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	message := c.Request().PostFormValue("message")
	if err := bot.SendMessage(playerID, message); err != nil {
		if err.Error() == "invalid parameters" {
			return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// GetFleetsHandler ...
func GetFleetsHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	fleets, _ := bot.GetFleets()
	return c.JSON(http.StatusOK, SuccessResp(fleets))
}

// GetSlotsHandler ...
func GetSlotsHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	slots := bot.GetSlots()
	return c.JSON(http.StatusOK, SuccessResp(slots))
}

// CancelFleetHandler ...
func CancelFleetHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	fleetID, err := strconv.ParseInt(c.Param("fleetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(bot.CancelFleet(FleetID(fleetID))))
}

// GetAttacksHandler ...
func GetAttacksHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	attacks, err := bot.GetAttacks()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(attacks))
}

// GalaxyInfosHandler ...
func GalaxyInfosHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	galaxy, err := strconv.ParseInt(c.Param("galaxy"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	system, err := strconv.ParseInt(c.Param("system"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	res, err := bot.GalaxyInfos(galaxy, system)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(res))
}

// GetResearchHandler ...
func GetResearchHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	return c.JSON(http.StatusOK, SuccessResp(bot.GetResearch()))
}

// BuyOfferOfTheDayHandler ...
func BuyOfferOfTheDayHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	if err := bot.BuyOfferOfTheDay(); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// GetMoonsHandler ...
func GetMoonsHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	return c.JSON(http.StatusOK, SuccessResp(bot.GetMoons()))
}

// GetMoonHandler ...
func GetMoonHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	moonID, err := strconv.ParseInt(c.Param("moonID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid moon id"))
	}
	moon, err := bot.GetMoon(moonID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid moon id"))
	}
	return c.JSON(http.StatusOK, SuccessResp(moon))
}

// GetMoonByCoordHandler ...
func GetMoonByCoordHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	galaxy, err := strconv.ParseInt(c.Param("galaxy"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid galaxy"))
	}
	system, err := strconv.ParseInt(c.Param("system"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid system"))
	}
	position, err := strconv.ParseInt(c.Param("position"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid position"))
	}
	planet, err := bot.GetMoon(Coordinate{Type: MoonType, Galaxy: galaxy, System: system, Position: position})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(planet))
}

// GetPlanetsHandler ...
func GetPlanetsHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	return c.JSON(http.StatusOK, SuccessResp(bot.GetPlanets()))
}

// GetCelestialItemsHandler ...
func GetCelestialItemsHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	celestialID, err := strconv.ParseInt(c.Param("celestialID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid celestial id"))
	}
	items, err := bot.GetItems(CelestialID(celestialID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(items))
}

// ActivateCelestialItemHandler ...
func ActivateCelestialItemHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	celestialID, err := strconv.ParseInt(c.Param("celestialID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid celestial id"))
	}
	ref := c.Param("itemRef")
	if err := bot.ActivateItem(ref, CelestialID(celestialID)); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// GetPlanetHandler ...
func GetPlanetHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	planet, err := bot.GetPlanet(PlanetID(planetID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(planet))
}

// GetPlanetByCoordHandler ...
func GetPlanetByCoordHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	galaxy, err := strconv.ParseInt(c.Param("galaxy"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid galaxy"))
	}
	system, err := strconv.ParseInt(c.Param("system"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid system"))
	}
	position, err := strconv.ParseInt(c.Param("position"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid position"))
	}
	planet, err := bot.GetPlanet(Coordinate{Type: PlanetType, Galaxy: galaxy, System: system, Position: position})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(planet))
}

// GetResourcesDetailsHandler ...
func GetResourcesDetailsHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	resources, err := bot.GetResourcesDetails(CelestialID(planetID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(resources))
}

// GetResourceSettingsHandler ...
func GetResourceSettingsHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	res, err := bot.GetResourceSettings(PlanetID(planetID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(res))
}

// SetResourceSettingsHandler ...
// curl 127.0.0.1:1234/bot/planets/123/resource-settings -d 'metalMine=100&crystalMine=100&deuteriumSynthesizer=100&solarPlant=100&fusionReactor=100&solarSatellite=100'
func SetResourceSettingsHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	metalMine, err := strconv.ParseInt(c.Request().PostFormValue("metalMine"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid metalMine"))
	}
	crystalMine, err := strconv.ParseInt(c.Request().PostFormValue("crystalMine"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid crystalMine"))
	}
	deuteriumSynthesizer, err := strconv.ParseInt(c.Request().PostFormValue("deuteriumSynthesizer"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid deuteriumSynthesizer"))
	}
	solarPlant, err := strconv.ParseInt(c.Request().PostFormValue("solarPlant"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid solarPlant"))
	}
	fusionReactor, err := strconv.ParseInt(c.Request().PostFormValue("fusionReactor"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid fusionReactor"))
	}
	solarSatellite, err := strconv.ParseInt(c.Request().PostFormValue("solarSatellite"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid solarSatellite"))
	}
	crawler, err := strconv.ParseInt(c.Request().PostFormValue("crawler"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid crawler"))
	}
	settings := ResourceSettings{
		MetalMine:            metalMine,
		CrystalMine:          crystalMine,
		DeuteriumSynthesizer: deuteriumSynthesizer,
		SolarPlant:           solarPlant,
		FusionReactor:        fusionReactor,
		SolarSatellite:       solarSatellite,
		Crawler:              crawler,
	}
	if err := bot.SetResourceSettings(PlanetID(planetID), settings); err != nil {
		if err == ErrInvalidPlanetID {
			return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// GetResourcesBuildingsHandler ...
func GetResourcesBuildingsHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	res, err := bot.GetResourcesBuildings(CelestialID(planetID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(res))
}

// GetDefenseHandler ...
func GetDefenseHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	res, err := bot.GetDefense(CelestialID(planetID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(res))
}

// GetShipsHandler ...
func GetShipsHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	res, err := bot.GetShips(CelestialID(planetID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(res))
}

// GetFacilitiesHandler ...
func GetFacilitiesHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	res, err := bot.GetFacilities(CelestialID(planetID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(res))
}

// BuildHandler ...
func BuildHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	ogameID, err := strconv.ParseInt(c.Param("ogameID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ogame id"))
	}
	nbr, err := strconv.ParseInt(c.Param("nbr"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid nbr"))
	}
	if err := bot.Build(CelestialID(planetID), ID(ogameID), nbr); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// BuildCancelableHandler ...
func BuildCancelableHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	ogameID, err := strconv.ParseInt(c.Param("ogameID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ogame id"))
	}
	if err := bot.BuildCancelable(CelestialID(planetID), ID(ogameID)); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// BuildProductionHandler ...
func BuildProductionHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	ogameID, err := strconv.ParseInt(c.Param("ogameID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ogame id"))
	}
	nbr, err := strconv.ParseInt(c.Param("nbr"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid nbr"))
	}
	if err := bot.BuildProduction(CelestialID(planetID), ID(ogameID), nbr); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// BuildBuildingHandler ...
func BuildBuildingHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	ogameID, err := strconv.ParseInt(c.Param("ogameID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ogame id"))
	}
	if err := bot.BuildBuilding(CelestialID(planetID), ID(ogameID)); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// BuildTechnologyHandler ...
func BuildTechnologyHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	ogameID, err := strconv.ParseInt(c.Param("ogameID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ogame id"))
	}
	if err := bot.BuildTechnology(CelestialID(planetID), ID(ogameID)); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// BuildDefenseHandler ...
func BuildDefenseHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	ogameID, err := strconv.ParseInt(c.Param("ogameID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ogame id"))
	}
	nbr, err := strconv.ParseInt(c.Param("nbr"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid nbr"))
	}
	if err := bot.BuildDefense(CelestialID(planetID), ID(ogameID), nbr); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// BuildShipsHandler ...
func BuildShipsHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	ogameID, err := strconv.ParseInt(c.Param("ogameID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ogame id"))
	}
	nbr, err := strconv.ParseInt(c.Param("nbr"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid nbr"))
	}
	if err := bot.BuildShips(CelestialID(planetID), ID(ogameID), nbr); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// GetProductionHandler ...
func GetProductionHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	res, _, err := bot.GetProduction(CelestialID(planetID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(res))
}

// ConstructionsBeingBuiltHandler ...
func ConstructionsBeingBuiltHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	buildingID, buildingCountdown, researchID, researchCountdown := bot.ConstructionsBeingBuilt(CelestialID(planetID))
	return c.JSON(http.StatusOK, SuccessResp(
		struct {
			BuildingID        int64
			BuildingCountdown int64
			ResearchID        int64
			ResearchCountdown int64
		}{
			BuildingID:        int64(buildingID),
			BuildingCountdown: buildingCountdown,
			ResearchID:        int64(researchID),
			ResearchCountdown: researchCountdown,
		},
	))
}

// CancelBuildingHandler ...
func CancelBuildingHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	if err := bot.CancelBuilding(CelestialID(planetID)); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// CancelResearchHandler ...
func CancelResearchHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	if err := bot.CancelResearch(CelestialID(planetID)); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// GetResourcesHandler ...
func GetResourcesHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	res, err := bot.GetResources(CelestialID(planetID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(res))
}

// GetPriceHandler ...
func GetPriceHandler(c echo.Context) error {
	ogameID, err := strconv.ParseInt(c.Param("ogameID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ogameID"))
	}
	nbr, err := strconv.ParseInt(c.Param("nbr"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid nbr"))
	}
	ogameObj := Objs.ByID(ID(ogameID))
	if ogameObj != nil {
		price := ogameObj.GetPrice(nbr)
		return c.JSON(http.StatusOK, SuccessResp(price))
	}
	return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ogameID"))
}

// SendFleetHandler ...
// curl 127.0.0.1:1234/bot/planets/123/send-fleet -d 'ships=203,1&ships=204,10&speed=10&galaxy=1&system=1&type=1&position=1&mission=3&metal=1&crystal=2&deuterium=3'
func SendFleetHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}

	if err := c.Request().ParseForm(); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid form"))
	}

	var ships []Quantifiable
	where := Coordinate{Type: PlanetType}
	mission := Transport
	var duration int64
	var unionID int64
	payload := Resources{}
	speed := HundredPercent
	for key, values := range c.Request().PostForm {
		switch key {
		case "ships":
			for _, s := range values {
				a := strings.Split(s, ",")
				shipID, err := strconv.ParseInt(a[0], 10, 64)
				if err != nil || !IsShipID(shipID) {
					return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ship id "+a[0]))
				}
				nbr, err := strconv.ParseInt(a[1], 10, 64)
				if err != nil || nbr < 0 {
					return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid nbr "+a[1]))
				}
				ships = append(ships, Quantifiable{ID: ID(shipID), Nbr: nbr})
			}
		case "speed":
			speedInt, err := strconv.ParseInt(values[0], 10, 64)
			if err != nil || speedInt < 0 || speedInt > 10 {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid speed"))
			}
			speed = Speed(speedInt)
		case "galaxy":
			galaxy, err := strconv.ParseInt(values[0], 10, 64)
			if err != nil {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid galaxy"))
			}
			where.Galaxy = galaxy
		case "system":
			system, err := strconv.ParseInt(values[0], 10, 64)
			if err != nil {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid system"))
			}
			where.System = system
		case "position":
			position, err := strconv.ParseInt(values[0], 10, 64)
			if err != nil {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid position"))
			}
			where.Position = position
		case "type":
			t, err := strconv.ParseInt(values[0], 10, 64)
			if err != nil {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid type"))
			}
			where.Type = CelestialType(t)
		case "mission":
			missionInt, err := strconv.ParseInt(values[0], 10, 64)
			if err != nil {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid mission"))
			}
			mission = MissionID(missionInt)
		case "duration":
			duration, err = strconv.ParseInt(values[0], 10, 64)
			if err != nil {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid duration"))
			}
		case "union":
			unionID, err = strconv.ParseInt(values[0], 10, 64)
			if err != nil {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid union id"))
			}
		case "metal":
			metal, err := strconv.ParseInt(values[0], 10, 64)
			if err != nil || metal < 0 {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid metal"))
			}
			payload.Metal = metal
		case "crystal":
			crystal, err := strconv.ParseInt(values[0], 10, 64)
			if err != nil || crystal < 0 {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid crystal"))
			}
			payload.Crystal = crystal
		case "deuterium":
			deuterium, err := strconv.ParseInt(values[0], 10, 64)
			if err != nil || deuterium < 0 {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid deuterium"))
			}
			payload.Deuterium = deuterium
		}
	}

	fleet, err := bot.SendFleet(CelestialID(planetID), ships, speed, where, mission, payload, duration, unionID)
	if err != nil &&
		(err == ErrInvalidPlanetID ||
			err == ErrNoShipSelected ||
			err == ErrUninhabitedPlanet ||
			err == ErrNoDebrisField ||
			err == ErrPlayerInVacationMode ||
			err == ErrAdminOrGM ||
			err == ErrNoAstrophysics ||
			err == ErrNoobProtection ||
			err == ErrPlayerTooStrong ||
			err == ErrNoMoonAvailable ||
			err == ErrNoRecyclerAvailable ||
			err == ErrNoEventsRunning ||
			err == ErrPlanetAlreadyReservedForRelocation) {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(fleet))
}

// GetAlliancePageContentHandler ...
func GetAlliancePageContentHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	allianceID := c.QueryParam("allianceId")
	vals := url.Values{"allianceId": {allianceID}}
	pageHTML, _ := bot.GetAlliancePageContent(vals)
	return c.HTML(http.StatusOK, string(pageHTML))
}

func replaceHostname(bot *OGame, html []byte) []byte {
	serverURLBytes := []byte(bot.serverURL)
	apiNewHostnameBytes := []byte(bot.apiNewHostname)
	escapedServerURL := bytes.Replace(serverURLBytes, []byte("/"), []byte(`\/`), -1)
	doubleEscapedServerURL := bytes.Replace(serverURLBytes, []byte("/"), []byte("\\\\\\/"), -1)
	escapedAPINewHostname := bytes.Replace(apiNewHostnameBytes, []byte("/"), []byte(`\/`), -1)
	doubleEscapedAPINewHostname := bytes.Replace(apiNewHostnameBytes, []byte("/"), []byte("\\\\\\/"), -1)
	html = bytes.Replace(html, serverURLBytes, apiNewHostnameBytes, -1)
	html = bytes.Replace(html, escapedServerURL, escapedAPINewHostname, -1)
	html = bytes.Replace(html, doubleEscapedServerURL, doubleEscapedAPINewHostname, -1)
	return html
}

// GetStaticHandler ...
func GetStaticHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)

	newURL := bot.serverURL + c.Request().URL.String()
	req, err := http.NewRequest("GET", newURL, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	resp, err := bot.Client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	defer resp.Body.Close()
	body, _, err := readBody(resp)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}

	// Copy the original HTTP headers to our client
	for k, vv := range resp.Header { // duplicate headers are acceptable in HTTP spec, so add all of them individually: https://stackoverflow.com/questions/4371328/are-duplicate-http-response-headers-acceptable
		k = http.CanonicalHeaderKey(k)
		if k != "Content-Length" && k != "Content-Encoding" { // https://github.com/alaingilbert/ogame/pull/80#issuecomment-674559853
			for _, v := range vv {
				c.Response().Header().Add(k, v)
			}
		}
	}

	if strings.Contains(c.Request().URL.String(), ".xml") {
		body = replaceHostname(bot, body)
		return c.Blob(http.StatusOK, "application/xml", body)
	}

	contentType := http.DetectContentType(body)
	if strings.Contains(newURL, ".css") {
		contentType = "text/css"
	} else if strings.Contains(newURL, ".js") {
		contentType = "application/javascript"
	} else if strings.Contains(newURL, ".gif") {
		contentType = "image/gif"
	}

	return c.Blob(http.StatusOK, contentType, body)
}

// GetFromGameHandler ...
func GetFromGameHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	vals := url.Values{"page": {"ingame"}, "component": {"overview"}}
	if len(c.QueryParams()) > 0 {
		vals = c.QueryParams()
	}
	pageHTML, _ := bot.GetPageContent(vals)
	pageHTML = replaceHostname(bot, pageHTML)
	return c.HTMLBlob(http.StatusOK, pageHTML)
}

// PostToGameHandler ...
func PostToGameHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	vals := url.Values{"page": {"ingame"}, "component": {"overview"}}
	if len(c.QueryParams()) > 0 {
		vals = c.QueryParams()
	}
	payload, _ := c.FormParams()
	pageHTML, _ := bot.PostPageContent(vals, payload)
	pageHTML = replaceHostname(bot, pageHTML)
	return c.HTMLBlob(http.StatusOK, pageHTML)
}

// GetStaticHEADHandler ...
func GetStaticHEADHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	newURL := "/api/" + strings.Join(c.ParamValues(), "") // + "?" + c.QueryString()
	if len(c.QueryString()) > 0 {
		newURL = newURL + "?" + c.QueryString()
	}
	headers, err := bot.HeadersForPage(newURL)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	if len(headers) < 1 {
		return c.NoContent(http.StatusFailedDependency)
	}
	// Copy the original HTTP HEAD headers to our client
	for k, vv := range headers { // duplicate headers are acceptable in HTTP spec, so add all of them individually: https://stackoverflow.com/questions/4371328/are-duplicate-http-response-headers-acceptable
		k = http.CanonicalHeaderKey(k)
		for _, v := range vv {
			c.Response().Header().Add(k, v)
		}
	}
	return c.NoContent(http.StatusOK)
}

// GetEmpireHandler ...
func GetEmpireHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	nbr, err := strconv.ParseInt(c.Param("typeID"), 10, 64)
	if err != nil || nbr > 1 {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid typeID"))
	}
	getEmpire, err := bot.GetEmpireJSON(nbr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(getEmpire))
}

// DeleteMessageHandler ...
func DeleteMessageHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	messageID, err := strconv.ParseInt(c.Param("messageID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid message id"))
	}
	if err := bot.DeleteMessage(messageID); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// DeleteEspionageMessagesHandler ...
func DeleteEspionageMessagesHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	if err := bot.DeleteAllMessagesFromTab(20); err != nil { // 20 = Espionage Reports
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "Unable to delete Espionage Reports"))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// DeleteMessagesFromTabHandler ...
func DeleteMessagesFromTabHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	tabIndex, err := strconv.ParseInt(c.Param("tabIndex"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "must provide tabIndex"))
	}
	if tabIndex < 20 || tabIndex > 24 {
		/*
			tabid: 20 => Espionage
			tabid: 21 => Combat Reports
			tabid: 22 => Expeditions
			tabid: 23 => Unions/Transport
			tabid: 24 => Other
		*/
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid tabIndex provided"))
	}
	if err := bot.DeleteAllMessagesFromTab(tabIndex); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "Unable to delete message from tab "+strconv.FormatInt(tabIndex, 10)))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// SendIPMHandler ...
func SendIPMHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	ipmAmount, err := strconv.ParseInt(c.Param("ipmAmount"), 10, 64)
	if err != nil || ipmAmount < 1 {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ipmAmount"))
	}
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil || planetID < 1 {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	galaxy, err := strconv.ParseInt(c.Request().PostFormValue("galaxy"), 10, 64)
	if err != nil || galaxy < 1 || galaxy > bot.serverData.Galaxies {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid galaxy"))
	}
	system, err := strconv.ParseInt(c.Request().PostFormValue("system"), 10, 64)
	if err != nil || system < 1 || system > bot.serverData.Systems {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid system"))
	}
	position, err := strconv.ParseInt(c.Request().PostFormValue("position"), 10, 64)
	if err != nil || position < 1 || position > 15 {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid position"))
	}
	planetTypeInt, err := strconv.ParseInt(c.Param("type"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	planetType := CelestialType(planetTypeInt)
	if planetType != PlanetType && planetType != MoonType { // only accept planet/moon types
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid type"))
	}
	priority, _ := strconv.ParseInt(c.Request().PostFormValue("priority"), 10, 64)
	coord := Coordinate{Type: planetType, Galaxy: galaxy, System: system, Position: position}
	duration, err := bot.SendIPM(PlanetID(planetID), coord, ipmAmount, ID(priority))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(duration))
}

// TeardownHandler ...
func TeardownHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil || planetID < 0 {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	ogameID, err := strconv.ParseInt(c.Param("ogameID"), 10, 64)
	if err != nil || planetID < 0 {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ogame id"))
	}
	if err = bot.TearDown(CelestialID(planetID), ID(ogameID)); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// GetAuctionHandler ...
func GetAuctionHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	auction, err := bot.GetAuction()
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "could not open auction page"))
	}
	return c.JSON(http.StatusOK, SuccessResp(auction))
}

// DoAuctionHandler (`celestialID=metal:crystal:deuterium` eg: `123456=123:456:789`)
func DoAuctionHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	bid := make(map[CelestialID]Resources)
	if err := c.Request().ParseForm(); err != nil { // Required for PostForm, not for PostFormValue
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid form"))
	}
	for key, values := range c.Request().PostForm {
		for _, s := range values {
			var metal, crystal, deuterium int64
			if n, err := fmt.Sscanf(s, "%d:%d:%d", &metal, &crystal, &deuterium); err != nil || n != 3 {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid bid format"))
			}
			celestialIDInt, err := strconv.ParseInt(key, 10, 64)
			if err != nil {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid celestial ID"))
			}
			bid[CelestialID(celestialIDInt)] = Resources{Metal: metal, Crystal: crystal, Deuterium: deuterium}
		}
	}
	if err := bot.DoAuction(bid); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// PhalanxHandler ...
func PhalanxHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	moonID, err := strconv.ParseInt(c.Param("moonID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid moon id"))
	}
	galaxy, err := strconv.ParseInt(c.Param("galaxy"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid galaxy"))
	}
	system, err := strconv.ParseInt(c.Param("system"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid system"))
	}
	position, err := strconv.ParseInt(c.Param("position"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid position"))
	}
	coord := Coordinate{Type: PlanetType, Galaxy: galaxy, System: system, Position: position}
	fleets, err := bot.Phalanx(MoonID(moonID), coord)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(fleets))
}

// JumpGateHandler ...
func JumpGateHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	if err := c.Request().ParseForm(); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid form"))
	}
	moonOriginID, err := strconv.ParseInt(c.Param("moonID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid origin moon id"))
	}
	moonDestinationID, err := strconv.ParseInt(c.Request().PostFormValue("moonDestination"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid destination moon id"))
	}
	var ships ShipsInfos
	for key, values := range c.Request().PostForm {
		switch key {
		case "ships":
			for _, s := range values {
				a := strings.Split(s, ",")
				shipID, err := strconv.ParseInt(a[0], 10, 64)
				if err != nil || !IsShipID(shipID) {
					return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ship id "+a[0]))
				}
				nbr, err := strconv.ParseInt(a[1], 10, 64)
				if err != nil || nbr < 0 {
					return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid nbr "+a[1]))
				}
				ships.Set(ID(shipID), nbr)
			}
		}
	}
	success, rechargeCountdown, err := bot.JumpGate(MoonID(moonOriginID), MoonID(moonDestinationID), ships)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(map[string]interface{}{
		"success":           success,
		"rechargeCountdown": rechargeCountdown,
	}))
}

// TechsHandler ...
func TechsHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	celestialID, err := strconv.ParseInt(c.Param("celestialID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid celestial id"))
	}
	supplies, facilities, ships, defenses, researches, err := bot.GetTechs(CelestialID(celestialID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(map[string]interface{}{
		"supplies":   supplies,
		"facilities": facilities,
		"ships":      ships,
		"defenses":   defenses,
		"researches": researches,
	}))
}

// GetCaptchaHandler ...
func GetCaptchaHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)

	gameEnvironmentID, platformGameID, err := getConfiguration(bot)
	if err != nil {
		return c.HTML(http.StatusOK, err.Error())
	}

	//var out postSessionsResponse
	payload := url.Values{
		"autoGameAccountCreation": {"false"},
		"gameEnvironmentId":       {gameEnvironmentID},
		"platformGameId":          {platformGameID},
		"gfLang":                  {"en"},
		"locale":                  {"en_GB"},
		"identity":                {bot.Username},
		"password":                {bot.password},
	}
	req, err := http.NewRequest("POST", "https://gameforge.com/api/v1/auth/thin/sessions", strings.NewReader(payload.Encode()))
	if err != nil {
		return c.HTML(http.StatusOK, err.Error())
	}

	if bot.otpSecret != "" {
		passcode, err := totp.GenerateCodeCustom(bot.otpSecret, time.Now(), totp.ValidateOpts{
			Period:    30,
			Skew:      1,
			Digits:    otp.DigitsSix,
			Algorithm: otp.AlgorithmSHA1,
		})
		if err != nil {
			return c.HTML(http.StatusOK, err.Error())
		}
		req.Header.Add("tnt-2fa-code", passcode)
		req.Header.Add("tnt-installation-id", "")
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")

	resp, err := bot.doReqWithLoginProxyTransport(req)
	if err != nil {
		return c.HTML(http.StatusOK, err.Error())
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	if resp.StatusCode == http.StatusForbidden {
		data403, _, _ := readBody(resp)
		return c.HTML(http.StatusOK, string(data403))
	}

	if resp.StatusCode == http.StatusConflict {
		var temp struct {
			ID          string `json:"id"`
			LastUpdated int    `json:"lastUpdated"`
			Status      string `json:"status"`
		}

		challengeID := resp.Header.Get(gfChallengeID)
		challengeID = strings.Replace(challengeID, ";https://challenge.gameforge.com", "", -1)

		req1, err := http.NewRequest("GET", "https://image-drop-challenge.gameforge.com/challenge/"+challengeID+"/en-GB", strings.NewReader(payload.Encode()))
		if err != nil {
			return c.HTML(http.StatusOK, err.Error())
		}
		resp1, err := bot.doReqWithLoginProxyTransport(req1)
		if err != nil {
			return c.HTML(http.StatusOK, err.Error())
		}
		defer resp1.Body.Close()

		data, _, _ := readBody(resp1)
		if err := json.Unmarshal(data, &temp); err != nil {
			return c.HTML(http.StatusOK, err.Error())
		}

		html := `<img style="background-color: black;" src="/bot/captcha/question/` + challengeID + `" /><br />
<img style="background-color: black;" src="/bot/captcha/icons/` + challengeID + `" /><br />
<form action="/bot/captcha/solve" method="POST">
	<input type="hidden" name="challenge_id" value="` + challengeID + `" />
	Enter 0,1,2 or 3 and press Enter <input type="number" name="answer" />" +
</form>` + challengeID

		return c.HTML(http.StatusOK, html)
	}
	return c.HTML(http.StatusOK, "no captcha found")
}

// GetCaptchaImgHandler ...
func GetCaptchaImgHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	challengeID := c.Param("challengeID")
	req, _ := http.NewRequest("GET", "https://image-drop-challenge.gameforge.com/challenge/"+challengeID+"/en-GB/drag-icons", nil)
	resp, _ := bot.doReqWithLoginProxyTransport(req)
	//IMG: https://image-drop-challenge.gameforge.com/challenge/9c5c46b2-e479-4f17-bd35-03bc4e5beefc/en-GB/drag-icons?1611748479816
	defer resp.Body.Close()
	data, _, _ := readBody(resp)
	if data == nil {
		return c.HTML(http.StatusNotFound, "File not Found")
	}
	return c.Blob(http.StatusOK, "image/png", data)
}

// GetCaptchaTextHandler ...
func GetCaptchaTextHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	challengeID := c.Param("challengeID")
	//TEXT: https://image-drop-challenge.gameforge.com/challenge/9c5c46b2-e479-4f17-bd35-03bc4e5beefc/en-GB/text?1611748479816
	req, _ := http.NewRequest("GET", "https://image-drop-challenge.gameforge.com/challenge/"+challengeID+"/en-GB/text", nil)
	resp, _ := bot.doReqWithLoginProxyTransport(req)
	defer resp.Body.Close()
	data, _, _ := readBody(resp)
	if data == nil {
		return c.HTML(http.StatusNotFound, "File not Found")
	}
	return c.Blob(http.StatusOK, "image/png", data)
}

// GetCaptchaSolverHandler ...
func GetCaptchaSolverHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	challengeID := c.Request().PostFormValue("challenge_id")
	answer := c.Request().PostFormValue("answer")
	payload := `{"answer":` + answer + `}`
	req, _ := http.NewRequest("POST", "https://image-drop-challenge.gameforge.com/challenge/"+challengeID+"/en-GB", strings.NewReader(payload))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	resp, _ := bot.doReqWithLoginProxyTransport(req)
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	if !bot.IsLoggedIn() {
		if err := bot.Login(); err != nil {
			bot.error(err)
		}
	}
	return c.Redirect(http.StatusTemporaryRedirect, "/")
}
