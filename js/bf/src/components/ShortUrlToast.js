import { Toast, ToastContainer } from 'react-bootstrap';
import { RiInformationLine, RiErrorWarningLine } from "react-icons/ri";

const ShortUrlToast = (props) => {
  return (
    <ToastContainer className="position-fixed bottom-0 end-0 p-4">
      <Toast className="text-white bg-primary border-0" onClose={() => props.hideInfo()} show={props.info !== ""} delay={3000} autohide>
        <Toast.Header >
          <RiInformationLine size={32} className="me-2" />
          <strong className="me-auto">Information</strong>
        </Toast.Header>
        <Toast.Body>{props.info}</Toast.Body>
      </Toast>
      <Toast className="text-white bg-danger border-0" onClose={() => props.hideError()} show={props.error !== ""} delay={3000} autohide>
        <Toast.Header >
          <RiErrorWarningLine size={32} className="me-2" />
          <strong className="me-auto">Error</strong>
        </Toast.Header>
        <Toast.Body>{props.error}</Toast.Body>
      </Toast>
    </ToastContainer>
  )
}

export default ShortUrlToast;