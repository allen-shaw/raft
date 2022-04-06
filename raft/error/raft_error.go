package error

import "strconv"

// RaftError Raft error code.
type RaftError int

const (
	// Unknown error
	UNKNOWN RaftError = -1

	// Success, no error.
	SUCCESS = 0

	// All Kinds of Timeout(Including Election_timeout, Timeout_now, Stepdown_timeout)
	ERAFTTIMEDOUT = 10001

	// Bad User State Machine
	ESTATEMACHINE = 10002

	// Catchup Failed
	ECATCHUP = 10003

	// Trigger step_down(Not All)
	ELEADERREMOVED = 10004

	// Leader Is Not In The New Configuration
	ESETPEER = 10005

	// Shut down
	ENODESHUTDOWN = 10006

	// Receive Higher Term Requests
	EHIGHERTERMREQUEST = 10007

	// Receive Higher Term Response
	EHIGHERTERMRESPONSE = 10008

	//Node Is In Error
	EBADNODE = 10009

	//
	//    /**
	//     * <pre>
	//     * Node Votes For Some Candidate
	//     * </pre>
	//     * <p>
	//     * <code>EVOTEFORCANDIDATE = 10010;</code>
	//     */
	//    EVOTEFORCANDIDATE(10010),
	//
	//    /**
	//     * <pre>
	//     * Follower(without leader) or Candidate Receives
	//     * </pre>
	//     * <p>
	//     * <code>ENEWLEADER = 10011;</code>
	//     */
	//    ENEWLEADER(10011),
	//
	//    /**
	//     * <pre>
	//     * Append_entries/Install_snapshot Request from a new leader
	//     * </pre>
	//     * <p>
	//     * <code>ELEADERCONFLICT = 10012;</code>
	//     */
	//    ELEADERCONFLICT(10012),
	//
	//    /**
	//     * <pre>
	//     * Trigger on_leader_stop
	//     * </pre>
	//     * <p>
	//     * <code>ETRANSFERLEADERSHIP = 10013;</code>
	//     */
	//    ETRANSFERLEADERSHIP(10013),
	//
	//    /**
	//     * <pre>
	//     * The log at the given index is deleted
	//     * </pre>
	//     * <p>
	//     * <code>ELOGDELETED = 10014;</code>
	//     */
	//    ELOGDELETED(10014),
	//
	//    /**
	//     * <pre>
	//     * No available user log to read
	//     * </pre>
	//     * <p>
	//     * <code>ENOMOREUSERLOG = 10015;</code>
	//     */
	//    ENOMOREUSERLOG(10015),
	//
	//    /* other non-raft error codes 1000~10000 */
	//    /**
	//     * Invalid rpc request
	//     */
	//    EREQUEST(1000),
	//
	//    /**
	//     * Task is stopped
	//     */
	//    ESTOP(1001),
	//
	//    /**
	//     * Retry again
	//     */
	//    EAGAIN(1002),
	//
	//    /**
	//     * Interrupted
	//     */
	//    EINTR(1003),
	//
	//    /**
	//     * Internal exception
	//     */
	//    EINTERNAL(1004),
	//
	//    /**
	//     * Task is canceled
	//     */
	//    ECANCELED(1005),
	//
	//    /**
	//     * Host is down
	//     */
	//    EHOSTDOWN(1006),
	//
	//    /**
	//     * Service is shutdown
	//     */
	//    ESHUTDOWN(1007),
	//
	//    /**
	//     * Permission issue
	//     */
	//    EPERM(1008),
	//
	//    /**
	//     * Server is in busy state
	//     */
	//    EBUSY(1009),
	//
	//    /**
	//     * Timed out
	//     */
	//    ETIMEDOUT(1010),
	//
	//    /**
	//     * Data is stale
	//     */
	//    ESTALE(1011),
	//
	//    /**
	//     * Something not found
	//     */
	//    ENOENT(1012),
	//
	//    /**
	//     * File/folder already exists
	//     */
	//    EEXISTS(1013),
	//
	//    /**
	//     * IO error
	//     */
	//    EIO(1014),
	//
	//    /**
	//     * Invalid value.
	//     */
	//    EINVAL(1015),
	//
	//    /**
	//     * Permission denied
	//     */
	//    EACCES(1016);
)

func DescribeCode(code int) string {
	return "<Unknown:" + strconv.Itoa(code) + ">"
}
