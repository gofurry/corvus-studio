import { Button, Space } from 'antd';
import { Link, useLocation } from 'react-router-dom';

type ProjectWorkspaceNavProps = {
  projectId: string;
};

function ProjectWorkspaceNav({ projectId }: ProjectWorkspaceNavProps) {
  const location = useLocation();
  const links = [
    { label: 'Overview', to: `/projects/${projectId}` },
    { label: 'Release', to: `/projects/${projectId}/release` },
    { label: 'Checklist', to: `/projects/${projectId}/checklist` },
  ];

  return (
    <nav aria-label="Project workspace">
      <Space wrap>
        {links.map((link) => (
          <Link key={link.to} to={link.to}>
            <Button
              type={location.pathname === link.to ? 'primary' : 'default'}
              aria-current={location.pathname === link.to ? 'page' : undefined}
            >
              {link.label}
            </Button>
          </Link>
        ))}
      </Space>
    </nav>
  );
}

export default ProjectWorkspaceNav;
