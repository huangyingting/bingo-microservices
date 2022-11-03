import { Section } from './Styled';

function Footer() {
  return (
    <Section color="#f4f4f4" className='text-center'>
      <h5 className="fw-bold">Copyright © {new Date().getFullYear()} Yingting Huang. All Rights Reserved</h5>
    </Section>
  )
}
export default Footer;