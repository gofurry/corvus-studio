import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Empty, Flex, List, Space, Spin, Tag, Typography } from 'antd';
import { Link } from 'react-router-dom';

import { projectsApi } from '../../api/projects';
import { projectStageLabels, projectStatusLabels } from './projectPresentation';
import { projectQueryKeys } from './projectQueries';

function ProjectListPage() {
  const projects = useQuery({
    queryKey: projectQueryKeys.all,
    queryFn: projectsApi.list,
  });

  return (
    <section className="page-stack" aria-labelledby="projects-heading">
      <Flex justify="space-between" align="flex-start" gap="middle" wrap>
        <div>
          <Typography.Text className="page-eyebrow">Project foundation</Typography.Text>
          <Typography.Title id="projects-heading" level={1}>
            Projects
          </Typography.Title>
          <Typography.Paragraph className="page-intro">
            Keep each game workspace visible without moving or changing its files.
          </Typography.Paragraph>
        </div>
        <Link to="/projects/new">
          <Button type="primary" size="large">
            Create project
          </Button>
        </Link>
      </Flex>

      {projects.isPending ? (
        <div className="center-state" role="status" aria-label="Loading projects">
          <Spin size="large" />
        </div>
      ) : projects.isError ? (
        <Alert
          type="error"
          showIcon
          title="Projects are unavailable"
          description={projects.error.message}
          action={
            <Button onClick={() => void projects.refetch()} size="small">
              Retry
            </Button>
          }
        />
      ) : projects.data.length === 0 ? (
        <Card className="empty-card">
          <Empty
            description={
              <Space orientation="vertical">
                <Typography.Text strong>No projects yet</Typography.Text>
                <Typography.Text type="secondary">
                  Reference an existing game directory to create your first workspace.
                </Typography.Text>
              </Space>
            }
          >
            <Link to="/projects/new">
              <Button type="primary">Create project</Button>
            </Link>
          </Empty>
        </Card>
      ) : (
        <List
          className="project-list"
          dataSource={projects.data}
          renderItem={(project) => (
            <List.Item>
              <Link className="project-card-link" to={`/projects/${project.id}`}>
                <Card className="project-card" hoverable>
                  <Flex justify="space-between" align="flex-start" gap="middle" wrap>
                    <Space orientation="vertical" size={4}>
                      <Typography.Title level={3}>{project.name}</Typography.Title>
                      <Typography.Text type="secondary">{project.location}</Typography.Text>
                    </Space>
                    <Space>
                      <Tag color="geekblue">{projectStageLabels[project.stage]}</Tag>
                      <Tag>{projectStatusLabels[project.status]}</Tag>
                    </Space>
                  </Flex>
                </Card>
              </Link>
            </List.Item>
          )}
        />
      )}
    </section>
  );
}

export default ProjectListPage;
