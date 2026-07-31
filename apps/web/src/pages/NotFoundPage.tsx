import { Button, Result } from 'antd';
import { Link } from 'react-router-dom';

function NotFoundPage() {
  return (
    <Result
      status="404"
      title="Page not found"
      subTitle="The requested Corvus Studio page does not exist."
      extra={
        <Link to="/projects">
          <Button type="primary">Back to projects</Button>
        </Link>
      }
    />
  );
}

export default NotFoundPage;
