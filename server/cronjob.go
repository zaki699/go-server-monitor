package server

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/decanus/array"
	"github.com/go-co-op/gocron"
	"visualon.com/go-server-monitor/config"
	"visualon.com/go-server-monitor/db"
	"visualon.com/go-server-monitor/models"
)

type Stat_aggregated struct {
	avg_fps         float64
	avg_drop_frames float64
	avg_dup_frames  float64
	avg_speed       float64
	session_id      int64
	avg_streams_qp  string
	created_at      time.Time
	ended_at        time.Time
}

func sliceAtoi(sa []string) ([]float64, error) {
	si := make([]float64, 0, len(sa))
	for _, a := range sa {
		i, err := strconv.ParseFloat(a, 64)
		if err != nil {
			return si, err
		}
		si = append(si, i)
	}
	return si, nil
}

func sliceFloat64ToSliceString(sa []float64) ([]string, error) {
	si := make([]string, 0, len(sa))
	for _, a := range sa {
		i := strconv.FormatFloat(a, 'g', -1, 64)
		si = append(si, i)
	}
	return si, nil
}

func numberOfStreamsQP(streamsQP string) (int, error) {
	if len(streamsQP) == 0 {
		return -1, fmt.Errorf("streamsQP is empty")
	}

	streams_qp := strings.Split(streamsQP, ",")

	if len(streams_qp) == 0 {
		return -1, fmt.Errorf("unable to parse streamsQP")
	}

	return len(streams_qp), nil
}

func splitStreamsQP(streamsQP string) ([]float64, error) {
	if len(streamsQP) == 0 {
		return nil, fmt.Errorf("streamsQP is empty")
	}

	streams_qp := strings.Split(streamsQP, ",")

	if len(streams_qp) == 0 {
		return nil, fmt.Errorf("unable to parse streamsQP")
	}

	streams_qp_int, err := sliceAtoi(streams_qp)

	if err != nil {
		return nil, err
	}

	return streams_qp_int, nil
}

func processStreamsQP(streamsQP []string) (string, error) {
	var b string

	if len(streamsQP) == 0 {
		return b, fmt.Errorf("streams QP is empty")
	}

	nbr_qp_streams, err := numberOfStreamsQP(streamsQP[0])

	if err != nil {
		return b, err
	}

	qpList := make([][]float64, 0)

	for _, element := range streamsQP {
		qp_int, err := splitStreamsQP(element)

		if err != nil {
			return b, err
		}

		if len(qp_int) > 0 {
			qpList = append(qpList, qp_int)
		}
	}

	qpList2 := make([][]float64, nbr_qp_streams, len(qpList))
	qpListAvg := make([]float64, nbr_qp_streams, len(qpList2))
	var nb = 0
	for nbr_qp_streams > 0 {
		if len(qpList) > 0 {
			for ib := range qpList {
				qpList2[nb] = append(qpList2[nb], qpList[ib][nb])
				continue
			}
		}

		qpListAvg[nb] = array.Average(qpList2[nb])
		nb++
		nbr_qp_streams--
	}

	streams_qp_s, err := sliceFloat64ToSliceString(qpListAvg)

	if err != nil {
		return b, err
	}

	streams_qp := strings.Join(streams_qp_s, ",")

	return streams_qp, nil
}

func AddStatAggregated(stat Stat_aggregated) (int64, error) {
	result, err := db.DB.Exec("INSERT INTO stats_aggregated (avg_fps, avg_drop_frames, avg_dup_frames, avg_speed, avg_streams_qp, session_id, created_at, ended_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", stat.avg_fps, stat.avg_drop_frames, stat.avg_dup_frames, stat.avg_speed, stat.avg_streams_qp, stat.session_id, stat.created_at, stat.ended_at)
	if err != nil {
		return 0, fmt.Errorf("AddStatAggregated: %v", err)
	}

	// Get the new stat's aggregated generated ID.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("AddStatAggregated: %v", err)
	}
	// Return the new stat aggregated ID.
	return id, nil
}

func DeleteStatRow(stat models.Stat) (int64, error) {
	result, err := db.DB.Exec("DELETE FROM stats where stat_id=?", stat.GetStatId())
	if err != nil {
		return 0, fmt.Errorf("DeleteStatRow: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("DeleteStatRow: %v", err)
	}
	// Return the deleted ID.
	return id, nil
}

func DeleteStatRowBySessionId(sessionId int64) (int64, error) {
	result, err := db.DB.Exec("DELETE FROM stats where session_id=?", sessionId)
	if err != nil {
		return 0, fmt.Errorf("DeleteStatRow: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("DeleteStatRow: %v", err)
	}
	// Return the deleted ID.
	return id, nil
}

func DeleteStatAggregatedRowBySessionId(sessionId int64) (int64, error) {
	result, err := db.DB.Exec("DELETE FROM stats_aggregated where session_id=?", sessionId)
	if err != nil {
		return 0, fmt.Errorf("DeleteStatRow: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("DeleteStatRow: %v", err)
	}
	// Return the deleted ID.
	return id, nil
}

func lastAggregatedStatRow(sessionId int64) (time.Time, error) {
	var created_at time.Time
	// Query for a value based on a single row.
	if err := db.DB.QueryRow("SELECT ended_at from stats_aggregated where session_id = ? ORDER BY stat_aggregated_id DESC LIMIT 1",
		sessionId).Scan(&created_at); err != nil {
		if err == sql.ErrNoRows {
			return time.Now(), fmt.Errorf("lastAggregatedStatRow %d: unknown aggregated stat", sessionId)
		}
		return time.Now(), fmt.Errorf("lastAggregatedStatRow %d", sessionId)
	}
	return created_at, nil
}

func lastAggregatedStatRowByDateFromAndTo(sessionId int64, from time.Time, to time.Time) bool {
	var created_at time.Time
	// Query for a value based on a single row.
	if err := db.DB.QueryRow("SELECT created_at from stats_aggregated where session_id = ? AND created_at = ? AND ended_at = ?",
		sessionId, from, to).Scan(&created_at); err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		return false
	}
	return true
}

func firstStatRow(sessionId int64) (time.Time, error) {
	var created_at time.Time
	// Query for a value based on a single row.
	if err := db.DB.QueryRow("SELECT created_at from stats where session_id = ? ORDER BY stat_id ASC LIMIT 1",
		sessionId).Scan(&created_at); err != nil {
		if err == sql.ErrNoRows {
			return time.Now(), fmt.Errorf("lastStatRow %d: unknown aggregated stat", sessionId)
		}
		return time.Now(), fmt.Errorf("lastStatRow %d", sessionId)
	}
	return created_at, nil
}

func lastStatRow(sessionId int64) (time.Time, error) {
	var created_at time.Time
	// Query for a value based on a single row.
	if err := db.DB.QueryRow("SELECT created_at from stats where session_id = ? ORDER BY stat_id DESC LIMIT 1",
		sessionId).Scan(&created_at); err != nil {
		if err == sql.ErrNoRows {
			return time.Now(), fmt.Errorf("lastStatRow %d: unknown aggregated stat", sessionId)
		}
		return time.Now(), fmt.Errorf("lastStatRow %d", sessionId)
	}
	return created_at, nil
}

func sessionByStatus(status string) ([]models.Session, error) {
	rows, err := db.DB.Query("SELECT session_id,created_at FROM sessions WHERE status = ?", status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// An album slice to hold data from returned rows.
	var sessions []models.Session

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var session models.Session
		if err := rows.Scan(&session.SessionId, &session.CreatedAt); err != nil {
			return sessions, err
		}
		sessions = append(sessions, session)
	}
	if err = rows.Err(); err != nil {
		return sessions, err
	}
	return sessions, nil
}

func statByDateInterval(from time.Time, to time.Time, sessionId int64) ([]models.Stat, error) {
	rows, err := db.DB.Query("SELECT stat_id,created_at,fps,speed,drop_frames,dup_frames,streams_qp,session_id FROM stats WHERE session_id=? AND created_at >= ? AND created_at < ?", sessionId, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// An stats slice to hold data from returned rows.
	var stats []models.Stat

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var stat models.Stat
		if err := rows.Scan(&stat.StatId, &stat.CreatedAt, &stat.Fps, &stat.Speed, &stat.DropFrames, &stat.DupFrames, &stat.StreamsQP, &stat.SessionId); err != nil {
			return stats, err
		}
		stats = append(stats, stat)
	}
	if err = rows.Err(); err != nil {
		return stats, err
	}
	return stats, nil
}

func deleteStatList(statList []models.Stat) error {
	if len(statList) > 0 {
		for _, element := range statList {
			_, err := DeleteStatRow(element)

			if err != nil {
				return err
			}
		}
	}
	return nil
}

func processStats(statList []models.Stat, sessionId int64, from time.Time, to time.Time) error {
	var fps []float64
	var drop_frames []int64
	var dup_frames []int64
	var speed []float64
	var qp []string

	if len(statList) == 0 {
		return fmt.Errorf("stats array is empty")
	}

	for _, element := range statList {
		fps = append(fps, element.GetFPS())
		drop_frames = append(drop_frames, element.GetDropFrames())
		dup_frames = append(dup_frames, element.GetDupFrames())
		speed = append(speed, element.GetSpeed())
		qp = append(qp, element.GetStreamsQP())
	}

	streams_qp, err := processStreamsQP(qp)

	if err != nil {
		return err
	}

	stat := Stat_aggregated{
		array.Average(fps),
		array.Average(drop_frames),
		array.Average(dup_frames),
		array.Average(speed),
		sessionId,
		streams_qp,
		from,
		to,
	}

	_, err = AddStatAggregated(stat)

	if err != nil {
		return err
	}
	return nil
}

func getStatsWithSession(session models.Session) error {
	var statList []models.Stat
	var dateTimeWhereWeStopAggregatingData time.Time
	var from time.Time
	var to time.Time

	interval := config.CONFIG.Cron.AggregatedInterval

	if interval <= 0 {
		interval = 60
	}

	timeDuration := time.Duration(interval)

	// Firstly we need to know at what date (time) we have to stop aggregating data
	// for that we have 2 ways :
	// Get the last row for the session n in  table stat aggregated
	// or Get the last row from the session n in table stats

	// Way 1
	dateTimeWhereWeStopAggregatingData, err := lastAggregatedStatRow(session.GetSessionId())

	if err != nil {
		fmt.Println("No date found in aggregated stats")
		// use the date from
		dateTimeWhereWeStopAggregatingData, err = firstStatRow(session.GetSessionId())
		if err != nil {
			return fmt.Errorf("could not find a date where to stop aggregating data")
		}
	}

	dateTimeWhere, err := lastStatRow(session.GetSessionId())

	if err != nil {
		return fmt.Errorf("could not find a date where to stop aggregating data")
	}

	fmt.Printf("Last row found in table stats_aggregated %s\n", dateTimeWhereWeStopAggregatingData)
	fmt.Printf("Last row found in table stats  %s\n", dateTimeWhere)
	// from is the createdAt date of the first row of stat for the session n.
	// This is our starting date.
	// Note that we will delete the rows in the table  stats once we agregrate the data.
	// it means the next run will have a date further.
	from = dateTimeWhereWeStopAggregatingData

	// to became from and to get one more hour
	to = from
	to = to.Add(time.Minute * timeDuration)

	for dateTimeWhere.Compare(to) > 0 {

		// check if the row for this interval does not exist already
		exist := lastAggregatedStatRowByDateFromAndTo(session.GetSessionId(), from, to)

		if exist {
			fmt.Printf("Data already exist for session Id : %d  from %s to %s \n", session.SessionId, from, to)
			from = to
			to = to.Add(time.Minute * timeDuration)
			continue
		}

		fmt.Printf("We are getting the data from %s to %s \n", from, to)

		statList, err = statByDateInterval(from, to, session.GetSessionId())

		if err != nil {
			fmt.Println(err)
			from = to
			to = to.Add(time.Minute * timeDuration)
			continue
		}

		err = processStats(statList, session.GetSessionId(), from, to)

		if err != nil {
			fmt.Println(err)
			from = to
			to = to.Add(time.Minute * timeDuration)
			continue
		}

		err = deleteStatList(statList)

		if err != nil {
			fmt.Print(err)
			return err
		}

		// from became to and to get one more hour
		from = to
		to = to.Add(time.Minute * timeDuration)

	}
	return nil
}

func getSessionList() error {

	fmt.Println("Starting Cron")
	var sessionList []models.Session
	sessionList, err := sessionByStatus("running")

	if err != nil {
		return err
	}

	for _, element := range sessionList {
		err := getStatsWithSession(element)

		if err != nil {
			fmt.Print(err)
			continue
		}
	}

	// clean DB
	sessionList, err = sessionByStatus("end")

	if err != nil {
		return err
	}

	if len(sessionList) > 0 {
		for _, element := range sessionList {
			_, err := DeleteStatRowBySessionId(element.GetSessionId())

			if err != nil {
				return err
			}
		}
	}
	return nil
}

func runCron() {
	// 3
	s := gocron.NewScheduler(time.UTC)

	// 4
	i := config.CONFIG.Cron.Interval

	if i <= 0 {
		i = 5
	}

	s.Every(i).Minutes().Do(func() {
		err := getSessionList()

		if err != nil {
			fmt.Print(err)
		}
	})

	// 5
	s.StartBlocking()

}
