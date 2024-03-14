package server

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"visualon.com/go-server-monitor/db"
	"visualon.com/go-server-monitor/models"
)

// Ingest progress
func Ingest(buf []string, sessionId int64) {

	if validBuf(buf) {
		s := models.NewStat()
		s.SessionId = sessionId
		ParseBuffer(buf, &s)
		err := createStat(s)

		if err != nil {
			fmt.Println(err)
		}

	} else {
		fmt.Printf("discard stat : %s\n", strings.Join(buf, "\n"))
	}
}

// Ingest progress
func IngestLog(buf string, sessionId int64) {

	if validLog(buf) {
		fmt.Printf("%+v [SessionId] %d | log : %+v\n", time_now(), sessionId, buf)
		l := models.NewLog()
		l.SessionId = int(sessionId)
		ParseLogBuffer(buf, &l)
		err := createLog(l)

		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Printf("%+v [SessionId] %d | discard log : %+v\n", time_now(), sessionId, buf)
	}
}

func CreateSession(url url.URL) (int64, error) {

	var err error

	session := models.New(url)

	err = checkSession(session)

	if err != nil {
		fmt.Println(err)
		return -1, err
	}

	q := "INSERT INTO `sessions` (channel_name, codec, definition, hostname, optimizer_enabled, status, created_at, preset, name, cmd) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"

	insert, err := db.DB.Prepare(q)
	if err != nil {
		fmt.Println(err)
		return -1, err
	}

	//convert string to bool
	var optimizer_enabled bool
	if n, err := strconv.ParseBool(session.GetOptimizerEnabled()); err == nil {
		optimizer_enabled = n
	}

	if err != nil {
		fmt.Println(err)
		optimizer_enabled = true
	}

	resp, err := insert.Exec(session.GetChannelName(), session.GetCodec(), session.GetDefinition(), session.GetHostName(), optimizer_enabled, session.GetStatus(), time_now(), session.GetPreset(), session.GetName(), session.GetCmd())

	insert.Close()
	if err != nil {
		fmt.Println(err)
		return -1, err
	}

	id, err := resp.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return -1, err
	}
	return id, nil
}

func EndSession(session_id int64) {
	var err error

	rows, err := db.DB.Query("select session_id,status from sessions where session_id=? and status = ?", session_id, "running")

	if err != nil {
		fmt.Println(err)
	}

	defer rows.Close()

	for rows.Next() {
		rows.Scan(&session_id, "running")
		go func(session_id int) {
			update, err := db.DB.Prepare("UPDATE sessions SET status=?, ended_at=? WHERE session_id=?")

			if err != nil {
				fmt.Println(err)
			}

			_, err = update.Exec("end", time_now(), session_id)

			if err != nil {
				fmt.Println(err)
			}

		}(int(session_id))
	}
}

func validBuf(buf []string) bool {
	if strings.HasPrefix(buf[0], "frame") {
		name := buf[len(buf)-1]
		return strings.HasSuffix(name, "continue")
	}
	return false
}

func validLog(buf string) bool {
	if len(buf) > 0 && (strings.Contains(buf, "[error]") || strings.Contains(buf, "[warning]")) {
		return !strings.Contains(buf, "x265 [info]")
	}
	return false
}

func ParseBuffer(buf []string, s *models.Stat) {
	log.Printf("[Session %d] Buffer Log : %s", s.SessionId, strings.Join([]string(buf), ","))

	res4 := buf
	var streamsQP []float64

	for _, a := range res4 {
		value := string(a[:])
		if strings.HasPrefix(value, "bitrate") {
			res4 := strings.Split(value, "=")
			if strings.Contains(res4[len(res4)-1], "N/A") {
				s.Bitrate = "0"
			} else {
				s.Bitrate = res4[len(res4)-1]
			}
		}

		if strings.HasPrefix(value, "frame") {
			res4 := strings.Split(value, "=")
			if n, err := strconv.ParseInt(res4[len(res4)-1], 10, 0); err == nil {
				s.Frames = n
			} else {
				s.Frames = 0
			}
		}

		if strings.HasPrefix(value, "dup_frame") {
			res4 := strings.Split(value, "=")
			if n, err := strconv.ParseInt(res4[len(res4)-1], 10, 0); err == nil {
				s.DupFrames = n
			} else {
				s.DupFrames = 0
			}
		}

		if strings.HasPrefix(value, "drop_frame") {
			res4 := strings.Split(value, "=")
			if n, err := strconv.ParseInt(res4[len(res4)-1], 10, 0); err == nil {
				s.DropFrames = n
			} else {
				s.DropFrames = 0
			}
		}

		if strings.HasPrefix(value, "fps") {
			res4 := strings.Split(value, "=")
			if n, err := strconv.ParseFloat(res4[len(res4)-1], 64); err == nil {
				s.Fps = n
			} else {
				fmt.Print(err)
				s.Fps = 0
			}
		}

		if strings.HasPrefix(value, "speed") {
			res4 := strings.Split(value, "=")
			res3 := strings.Split(res4[len(res4)-1], "x")
			if n, err := strconv.ParseFloat(strings.Trim(res3[len(res4)-2], " "), 64); err == nil {
				s.Speed = n
			} else {
				s.Speed = 0.0
			}
		}

		if strings.HasPrefix(value, "out_time_ms") {
			res4 := strings.Split(value, "=")
			if n, err := strconv.ParseInt(res4[len(res4)-1], 10, 0); err == nil {
				s.EncodingTime = n
			}
		}

		if strings.HasPrefix(value, "stream") {
			res4 := strings.Split(value, "=")
			if n, err := strconv.ParseFloat(res4[len(res4)-1], 64); err == nil {
				streamsQP = append(streamsQP, n)
			}
		}
	}

	s.StreamsQP = JoinStreamsQP(streamsQP)
}

func JoinStreamsQP(streamsQP []float64) string {
	valuesText := []string{}
	for i := range streamsQP {
		number := streamsQP[i]
		text := fmt.Sprintf("%f", number)
		valuesText = append(valuesText, text)
	}

	return strings.Join(valuesText, ",")
}

func createStat(s models.Stat) error {

	var err error

	err = checkStat(s)

	if err != nil {
		fmt.Println(err)
		return err
	}

	q := "INSERT INTO `stats` (frames, drop_frames, dup_frames, session_id, speed, bitrate, encoding_time, streams_qp, fps, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"
	insert, err := db.DB.Prepare(q)
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = insert.Exec(s.GetFrames(), s.GetDropFrames(), s.GetDupFrames(), s.GetSessionId(), s.GetSpeed(), s.GetBitrate(), s.GetEncodingTime(), s.GetStreamsQP(), s.GetFPS(), time_now())
	insert.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func time_now() string {
	t := time.Now()
	return t.Format("2006-01-02 15:04:05")
}

func CheckifSessionAlreadyExist(url url.URL) (int64, error) {

	var err error
	var session_idn int64 = -1

	var name = url.Query().Get("name")

	rows, err := db.DB.Query("select session_id from sessions where name=? and status=? limit 1", name, "running")

	if err != nil {
		fmt.Println(err)
		return session_idn, err
	}

	defer rows.Close()

	for rows.Next() {
		var session_id int64
		if err := rows.Scan(&session_id); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			log.Fatal(err)
		}
		session_idn = session_id
	}

	// If the database is being written to ensure to check for Close
	// errors that may be returned from the driver. The query may
	// encounter an auto-commit error and be forced to rollback changes.
	rerr := rows.Close()
	if rerr != nil {
		log.Fatal(rerr)
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return session_idn, nil
}

func ParseLogBuffer(buf string, l *models.Log) {

	pattern := "\\[(.*?)\\]\\s+"

	regex := regexp.MustCompile(pattern)
	result := regex.Split(buf, -1)

	l.Message = result[len(result)-1]

	pattern2 := "\\[([^][]*)]"

	// Probbly needs to be reviewed :)
	if n, err := regexp.Compile(pattern2); err == nil {
		_find := n.FindAllStringSubmatch(buf, -1)
		if len(_find) == 2 {
			l.Level = _find[len(_find)-1][len(_find[len(_find)-1])-1]
			var module = _find[0][len(_find[len(_find)-2])-1]
			var split_module = strings.Split(module, "@")
			l.Module = split_module[0]
		}

		if len(_find) == 1 {
			l.Level = _find[len(_find)-1][len(_find[len(_find)-1])-1]
		}
	}

}

func createLog(l models.Log) error {

	var err error

	err = checkLog(l)

	if err != nil {
		fmt.Println(err)
		return err
	}

	q := "INSERT INTO `logs` (level, message, module, session_id, created_at) VALUES (?, ?, ?, ?, ?);"
	insert, err := db.DB.Prepare(q)
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = insert.Exec(l.GetLevel(), l.GetMessage(), l.GetModule(), l.GetSessionId(), time_now())
	insert.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func checkSession(s models.Session) error {

	if len(s.ChannelName) == 0 {
		return errors.New("url param channel_name not found")
	}

	if len(s.Codec) == 0 {
		return errors.New("url param codec not found")
	}

	if len(s.Definition) == 0 {
		return errors.New("url param definition not found")
	}

	if len(s.HostName) == 0 {
		return errors.New("url param hostname not found")
	}

	if len(s.Name) == 0 {
		return errors.New("url param name not found")
	}

	if len(s.Preset) == 0 {
		return errors.New("url param preset not found")
	}

	if len(s.OptimizerEnabled) == 0 {
		return errors.New("url param optimizer_enabled not found")
	}

	if len(s.Status) == 0 {
		return errors.New("url param status not found")
	}

	if len(s.Cmd) == 0 {
		return errors.New("url param cmd not found")
	}

	return nil
}

func checkStat(s models.Stat) error {

	if s.Frames == -1 {
		return errors.New("frame not found")
	}

	if s.DropFrames == -1 {
		return errors.New("drop frame not found")
	}

	if s.DupFrames == -1 {
		return errors.New("duplicate frame not found")
	}

	if s.Speed == -1 {
		return errors.New("speed not found")
	}

	if s.Fps == -1 {
		return errors.New("fps not found")
	}

	if len(s.Bitrate) == 0 {
		return errors.New("bitrate not found")
	}

	if s.SessionId == -1 {
		return errors.New("sessionId not found")
	}

	if s.EncodingTime == -1 {
		return errors.New("EncodingTime not found")
	}

	if len(s.StreamsQP) == 0 {
		return errors.New("StreamsQP not found")
	}

	return nil
}

func checkLog(l models.Log) error {

	if len(l.Level) == 0 {
		return errors.New("level not found")
	}

	if len(l.Message) == 0 {
		return errors.New("message not found")
	}

	if len(l.Module) == 0 {
		return errors.New("module not found")
	}

	if l.SessionId == -1 {
		return errors.New("session not found")
	}

	return nil
}
