import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  Alert,
  Button,
  Card,
  Col,
  Flex,
  List,
  Popconfirm,
  Progress,
  Row,
  Skeleton,
  Space,
  Statistic,
  Tag,
  Typography,
} from 'antd';
import { Link, useNavigate, useParams } from 'react-router-dom';

import { projectsApi } from '../../api/projects';
import { releasesApi, type ReleaseGoal, type ReleaseStatus } from '../../api/releases';
import ProjectWorkspaceNav from '../projects/ProjectWorkspaceNav';
import { projectQueryKeys } from '../projects/projectQueries';
import {
  releaseStatusLabels,
  releaseTransitions,
  requirementLevelLabels,
  sourceLabels,
} from './workflowPresentation';
import { checklistQueryKeys, releaseQueryKeys } from './workflowQueries';

const templateKey = 'steam-coming-soon';

function ReleasePage() {
  const { projectId = '' } = useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const project = useQuery({
    queryKey: projectQueryKeys.detail(projectId),
    queryFn: () => projectsApi.get(projectId),
    enabled: projectId !== '',
  });
  const releases = useQuery({
    queryKey: releaseQueryKeys.project(projectId),
    queryFn: () => releasesApi.listForProject(projectId),
    enabled: projectId !== '',
  });
  const release = releases.data?.find((candidate) => candidate.goal_type === 'steam_coming_soon');
  const template = useQuery({
    queryKey: releaseQueryKeys.template(templateKey),
    queryFn: () => releasesApi.getTemplate(templateKey),
    enabled: releases.isSuccess && release === undefined,
  });

  const createRelease = useMutation({
    mutationFn: () => {
      if (template.data === undefined) {
        throw new Error('Release template is not loaded');
      }
      return releasesApi.create({
        project_id: projectId,
        goal_type: 'steam_coming_soon',
        template_key: template.data.key,
        template_version: template.data.version,
      });
    },
    onSuccess: (created) => {
      queryClient.setQueryData<ReleaseGoal[]>(releaseQueryKeys.project(projectId), [created]);
      queryClient.setQueryData(releaseQueryKeys.detail(created.id), created);
      void navigate(`/projects/${projectId}/checklist`);
    },
  });

  const transition = useMutation({
    mutationFn: ({ releaseId, status }: { releaseId: string; status: ReleaseStatus }) =>
      releasesApi.transition(releaseId, status),
    onSuccess: (updated) => {
      queryClient.setQueryData(releaseQueryKeys.detail(updated.id), updated);
      queryClient.setQueryData<ReleaseGoal[]>(releaseQueryKeys.project(projectId), (current) =>
        (current ?? []).map((candidate) => (candidate.id === updated.id ? updated : candidate)),
      );
      void queryClient.invalidateQueries({
        queryKey: checklistQueryKeys.list(updated.id, {}),
      });
    },
  });

  if (project.isPending || releases.isPending || (release === undefined && template.isPending)) {
    return <LoadingWorkspace />;
  }

  const error = project.error ?? releases.error ?? (release === undefined ? template.error : null);
  if (error !== null) {
    return (
      <section className="page-stack">
        <Alert
          type="error"
          showIcon
          title="Release workspace is unavailable"
          description={error.message}
          action={<Button onClick={() => window.location.reload()}>Retry</Button>}
        />
      </section>
    );
  }

  if (project.data === undefined) {
    return null;
  }

  return (
    <section className="page-stack" aria-labelledby="release-heading">
      <Flex justify="space-between" align="flex-start" gap="middle" wrap>
        <div>
          <Typography.Text className="page-eyebrow">{project.data.name}</Typography.Text>
          <Typography.Title id="release-heading" level={1}>
            Steam Coming Soon
          </Typography.Title>
          <Typography.Paragraph className="page-intro">
            Prepare the storefront work as one durable Release Goal and Checklist.
          </Typography.Paragraph>
        </div>
        <Link to="/projects">
          <Button>All projects</Button>
        </Link>
      </Flex>

      <ProjectWorkspaceNav projectId={projectId} />

      {release === undefined && template.data !== undefined ? (
        <TemplatePreview
          template={template.data}
          steamAppId={project.data.steam_app_id}
          creating={createRelease.isPending}
          error={createRelease.error}
          onCreate={() => createRelease.mutate()}
        />
      ) : release !== undefined ? (
        <ReleaseWorkspace
          release={release}
          transitionError={transition.error}
          transitioning={transition.isPending}
          onTransition={(status) => transition.mutate({ releaseId: release.id, status })}
          projectId={projectId}
        />
      ) : null}
    </section>
  );
}

function LoadingWorkspace() {
  return (
    <section className="page-stack" aria-label="Loading release workspace">
      <Card className="detail-card">
        <Skeleton active />
      </Card>
    </section>
  );
}

type TemplatePreviewProps = {
  template: Awaited<ReturnType<typeof releasesApi.getTemplate>>;
  steamAppId?: number | null;
  creating: boolean;
  error: Error | null;
  onCreate: () => void;
};

function TemplatePreview({
  template,
  steamAppId,
  creating,
  error,
  onCreate,
}: TemplatePreviewProps) {
  const steamCount = template.items.filter((item) => item.source === 'platform_template').length;
  const corvusCount = template.items.filter((item) => item.source === 'corvus_template').length;

  return (
    <>
      {steamAppId == null && (
        <Alert
          type="info"
          showIcon
          title="Steam App ID is not set"
          description="You can still create this workspace. Its Setup tasks will track Steamworks App preparation."
        />
      )}
      {error !== null && (
        <Alert
          type="error"
          showIcon
          title="Release workspace was not created"
          description={error.message}
        />
      )}
      <Card className="detail-card" title={template.name}>
        <Space orientation="vertical" size="large" className="full-width">
          <Typography.Paragraph>{template.description}</Typography.Paragraph>
          <Row gutter={[16, 16]}>
            <Col xs={12} md={6}>
              <Statistic title="Template version" value={template.version} />
            </Col>
            <Col xs={12} md={6}>
              <Statistic title="Total tasks" value={template.items.length} />
            </Col>
            <Col xs={12} md={6}>
              <Statistic title="Steam tasks" value={steamCount} />
            </Col>
            <Col xs={12} md={6}>
              <Statistic title="Corvus tasks" value={corvusCount} />
            </Col>
          </Row>
          <Typography.Text type="secondary">
            Reviewed {template.reviewed_at}. The template is copied into this Project and will not
            change automatically.
          </Typography.Text>
          <List
            className="template-task-list"
            dataSource={template.items}
            renderItem={(item) => (
              <List.Item>
                <List.Item.Meta
                  title={item.title}
                  description={
                    <Space wrap>
                      <Tag>{sourceLabels[item.source]}</Tag>
                      <Tag color={item.requirement_level === 'required' ? 'gold' : 'default'}>
                        {requirementLevelLabels[item.requirement_level]}
                      </Tag>
                    </Space>
                  }
                />
              </List.Item>
            )}
          />
          <Popconfirm
            title="Generate this release workspace?"
            description="This creates one Release Goal and all 12 Checklist tasks in a single operation."
            okText="Generate workspace"
            onConfirm={onCreate}
          >
            <Button type="primary" size="large" loading={creating}>
              Generate release workspace
            </Button>
          </Popconfirm>
        </Space>
      </Card>
    </>
  );
}

type ReleaseWorkspaceProps = {
  release: ReleaseGoal;
  transitionError: Error | null;
  transitioning: boolean;
  onTransition: (status: ReleaseStatus) => void;
  projectId: string;
};

function ReleaseWorkspace({
  release,
  transitionError,
  transitioning,
  onTransition,
  projectId,
}: ReleaseWorkspaceProps) {
  const summary = release.checklist_summary;
  const percent = summary.total === 0 ? 0 : Math.round((summary.done / summary.total) * 100);

  return (
    <>
      {transitionError !== null && (
        <Alert
          type="warning"
          showIcon
          title="Release status was not changed"
          description={transitionError.message}
        />
      )}
      <Card className="detail-card" title="Release status">
        <Space orientation="vertical" size="large" className="full-width">
          <Flex justify="space-between" align="center" gap="middle" wrap>
            <div>
              <Tag color="geekblue">{releaseStatusLabels[release.status]}</Tag>
              <Typography.Text type="secondary">
                Template {release.template_key} v{release.template_version}
              </Typography.Text>
            </div>
            <Link to={`/projects/${projectId}/checklist`}>
              <Button type="primary">Open checklist</Button>
            </Link>
          </Flex>
          <Progress percent={percent} status={summary.blocked > 0 ? 'exception' : 'active'} />
          <Row gutter={[16, 16]}>
            <Col xs={12} md={6}>
              <Statistic title="Completed" value={`${summary.done}/${summary.total}`} />
            </Col>
            <Col xs={12} md={6}>
              <Statistic
                title="Required"
                value={`${summary.required_done}/${summary.required_total}`}
              />
            </Col>
            <Col xs={12} md={6}>
              <Statistic title="Blocked" value={summary.blocked} />
            </Col>
          </Row>
          {release.status === 'preparing' && summary.required_done < summary.required_total && (
            <Alert
              type="info"
              showIcon
              title="Required tasks are still open"
              description="Complete every Required task before moving this Release Goal to Ready for review."
            />
          )}
          <Space wrap aria-label="Release status actions">
            {releaseTransitions[release.status].map((status) => (
              <Button key={status} loading={transitioning} onClick={() => onTransition(status)}>
                Move to {releaseStatusLabels[status]}
              </Button>
            ))}
          </Space>
        </Space>
      </Card>
    </>
  );
}

export default ReleasePage;
