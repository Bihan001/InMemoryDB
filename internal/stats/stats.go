package stats

import "github.com/Bihan001/MyDB/internal/interfaces"

type defaultStatsManager struct {
    DatabaseStats [4]map[string]int
}

func GetNewStatsManager() interfaces.StatsManager {
    statsManager := &defaultStatsManager{
        DatabaseStats: [4]map[string]int{},
    }
    statsManager.DatabaseStats[0] = make(map[string]int)
    return statsManager
}

func (sh *defaultStatsManager) RecordDBStat(metric string, value int) {
    sh.DatabaseStats[0][metric] = value
}

func (sh *defaultStatsManager) IncrDBStat(metric string) {
    sh.DatabaseStats[0][metric]++
}

func (sh *defaultStatsManager) DecrDBStat(metric string) {
    sh.DatabaseStats[0][metric]--
}

func (sh *defaultStatsManager) GetDbStats() [4]map[string]int {
    return sh.DatabaseStats
}
