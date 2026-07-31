import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Descriptions, Flex, Skeleton, Space, Tag, Typography } from 'antd';
import { Link, useParams } from 'react-router-dom';

import { projectsApi } from '../../api/projects';
import { formatProjectDate, projectStageLabels, projectStatusLabels } from './projectPresentation';
import { projectQueryKeys } from './projectQueries';

function ProjectDetailPage() {
  const { projectId = '' } = useParams();
  const project = useQuery({
    queryKey: projectQueryKeys.detail(projectId),
    queryFn: () => projectsApi.get(projectId),
    enabled: projectId !== '',
    staleTime: 30_000,
  });

  if (project.isPending) {
    return (
      <section className="page-stack narrow-page" aria-label="Loading project">
        <Card>
          <Skeleton active />
        </Card>
      </section>
    );
  }

  if (project.isError) {
    return (
      <section className="page-stack narrow-page">
        <Alert
          type="error"
          showIcon
          title="Project could not be opened"
          description={project.error.message}
          action={
            <Space>
              <Button onClick={() => void project.refetch()} size="small">
                Retry
              </Button>
              <Link to="/projects">
                <Button size="small">Back to projects</Button>
              </Link>
            </Space>
          }
        />
      </section>
    );
  }

  return (
    <section className="page-stack narrow-page" aria-labelledby="project-heading">
      <Flex justify="space-between" align="flex-start" gap="middle" wrap>
        <div>
          <Typography.Text className="page-eyebrow">Project details</Typography.Text>
          <Typography.Title id="project-heading" level={1}>
            {project.data.name}
          </Typography.Title>
          <Space>
            <Tag color="geekblue">{projectStageLabels[project.data.stage]}</Tag>
            <Tag>{projectStatusLabels[project.data.status]}</Tag>
          </Space>
        </div>
        <Link to="/projects">
          <Button>Back to projects</Button>
        </Link>
      </Flex>

      <Card className="detail-card">
        <Descriptions column={1} bordered>
          <Descriptions.Item label="Description">
            {project.data.description || 'No description'}
          </Descriptions.Item>
          <Descriptions.Item label="Project location">
            <Typography.Text code copyable>
              {project.data.location}
            </Typography.Text>
          </Descriptions.Item>
          <Descriptions.Item label="Primary language">{project.data.language}</Descriptions.Item>
          <Descriptions.Item label="Steam AppID">
            {project.data.steam_app_id ?? 'Not set'}
          </Descriptions.Item>
          <Descriptions.Item label="Created">
            {formatProjectDate(project.data.created_at)}
          </Descriptions.Item>
          <Descriptions.Item label="Last updated">
            {formatProjectDate(project.data.updated_at)}
          </Descriptions.Item>
        </Descriptions>
      </Card>
    </section>
  );
}

export default ProjectDetailPage;
