insert into roles ( id, name, description, created_at, updated_at ) values ( '2ea4415c-9748-493f-91ba-4a64506b7be8', 'tks-admin', 'tks-admin', now(), now() );
insert into roles ( id, name, description, created_at, updated_at ) values ( 'b2b689f0-ceeb-46c2-b280-0bc06896acd1', 'admin', 'admin', now(), now() );
insert into roles ( id, name, description, created_at, updated_at ) values ( 'd3015140-2b12-487a-9516-cdeed7c17735', 'project-leader', 'project-leader', now(), now() );
insert into roles ( id, name, description, created_at, updated_at ) values ( 'f6637d3d-3a0e-4db0-9086-c1b6dc9d433d', 'project-member', 'project-member', now(), now() );
insert into roles ( id, name, description, created_at, updated_at ) values ( 'b7ac7e7d-d8bc-470d-b6b2-3e0cc8ba55cc', 'project-viewer', 'project-viewer', now(), now() );
insert into roles ( id, name, description, created_at, updated_at ) values ( 'ff4187a2-f3c1-46b3-8448-03a4b5e132e7', 'user', 'user', now(), now() );

insert into policies ( role_id, name, description, c, create_priviledge, u, update_priviledge, r, read_priviledge, d, delete_priviledge, creator, created_at, updated_at ) values ( '2ea4415c-9748-493f-91ba-4a64506b7be8', 'organization', 'organization', 't', '', 't', '', 't', '', 't', '', '', now(), now() );
insert into policies ( role_id, name, description, c, create_priviledge, u, update_priviledge, r, read_priviledge, d, delete_priviledge, creator, created_at, updated_at ) values ( '2ea4415c-9748-493f-91ba-4a64506b7be8', 'project', 'project', 't', '', 't', '', 't', '', 't', '', '', now(), now() );
insert into policies ( role_id, name, description, c, create_priviledge, u, update_priviledge, r, read_priviledge, d, delete_priviledge, creator, created_at, updated_at ) values ( '2ea4415c-9748-493f-91ba-4a64506b7be8', 'user', 'user', 't', '', 't', '', 't', '', 't', '', '', now(), now() );
insert into policies ( role_id, name, description, c, create_priviledge, u, update_priviledge, r, read_priviledge, d, delete_priviledge, creator, created_at, updated_at ) values ( '2ea4415c-9748-493f-91ba-4a64506b7be8', 'cluster', 'cluster', 't', '', 't', '', 't', '', 't', '', '', now(), now() );
insert into policies ( role_id, name, description, c, create_priviledge, u, update_priviledge, r, read_priviledge, d, delete_priviledge, creator, created_at, updated_at ) values ( '2ea4415c-9748-493f-91ba-4a64506b7be8', 'service', 'service', 't', '', 't', '', 't', '', 't', '', '', now(), now() );
insert into policies ( role_id, name, description, c, create_priviledge, u, update_priviledge, r, read_priviledge, d, delete_priviledge, creator, created_at, updated_at ) values ( '2ea4415c-9748-493f-91ba-4a64506b7be8', 'k8s_resources', 'resources of k8s', 'f', '', 'f', '', 'f', '', 'f', '', '', now(), now() );

insert into organizations ( id, name, description, created_at, updated_at ) values ( 'master', 'master', 'tks', now(), now() );
insert into users ( id, account_id, name, password, organization_id, role_id, created_at, updated_at  ) values ( 'bf67de40-ce15-4dc0-b6c2-17f053ca504f', 'admin', 'admin', '$2a$10$Akf03nbLHk93sTtozm35XuINXkJeNX7A1T9o/Pxpg9R2B2PToBPOO', 'master', 'b2b689f0-ceeb-46c2-b280-0bc06896acd1', now(), now() );

insert into cloud_accounts ( id, name, description, organization_id, cloud_service, resource, created_at, updated_at ) 
values ( 'ce9e0387-01cb-4f37-a22a-fb91b6338434', 'aws', 'aws_description', 'master', 'AWS', 'result', now(), now() );

insert into stack_templates ( id, organization_id, name, description, version, cloud_service, platform, template, kube_version, kube_type, created_at, updated_at, services )
values ( '49901092-be76-4d4f-94e9-b84525f560b5', 'master', 'AWS Standard (x86)', 'included LMA', 'v1', 'AWS', 'x86', 'aws-reference', 'v1.24', 'AWS', now(), now(), '[{"name": "Logging,Monitoring,Alerting", "type": "LMA", "applications": [{"name": "prometheus-stack", "version": "v.44.3.1", "description": "통계데이터 제공을 위한 backend  플랫폼"}, {"name": "elastic-system", "version": "v1.8.0", "description": "로그 데이터 적재를 위한 Storage"}, {"name": "alertmanager", "version": "v0.23.0", "description": "Alert 관리를 위한 backend 서비스"}, {"name": "grafana", "version": "v6.50.7", "description": "모니터링 통합 포탈"}]}]' );
insert into stack_templates ( id, organization_id, name, description, version, cloud_service, platform, template, kube_version, kube_type, created_at, updated_at, services )
values ( '44d5e76b-63db-4dd0-a16e-11bd3f6054cf', 'master', 'AWS MSA Standard (x86)', 'included LMA, SERVICE MESH', 'v1', 'AWS', 'x86', 'aws-msa-reference', 'v1.24', 'AWS', now(), now(), '[{"name": "Logging,Monitoring,Alerting", "type": "LMA", "applications": [{"name": "prometheus-stack", "version": "v.44.3.1", "description": "통계데이터 제공을 위한 backend  플랫폼"}, {"name": "elastic-system", "version": "v1.8.0", "description": "로그 데이터 적재를 위한 Storage"}, {"name": "alertmanager", "version": "v0.23.0", "description": "Alert 관리를 위한 backend 서비스"}, {"name": "grafana", "version": "v6.50.7", "description": "모니터링 통합 포탈"}]}, {"name": "MSA", "type": "SERVICE_MESH", "applications": [{"name": "istio", "version": "v1.13.1", "description": "MSA 플랫폼"}, {"name": "jagger", "version": "v2.27.1", "description": "분산 서비스간 트랜잭션 추적을 위한 로깅 플랫폼"}, {"name": "kiali", "version": "v1.45.1", "description": "MSA 통합 모니터링포탈"}]}]' );
insert into stack_templates ( id, organization_id, name, description, version, cloud_service, platform, template, kube_version, kube_type, created_at, updated_at, services )
values ( 'fe1d97e0-7428-4be6-9c69-310a88b4ff46', 'master', 'AWS Standard (arm)', 'included LMA', 'v2', 'AWS', 'arm', 'aws-arm-reference', 'v1.24', 'EKS', now(), now(), '[{"name": "Logging,Monitoring,Alerting", "type": "LMA", "applications": [{"name": "prometheus-stack", "version": "v.44.3.1", "description": "통계데이터 제공을 위한 backend  플랫폼"}, {"name": "elastic-system", "version": "v1.8.0", "description": "로그 데이터 적재를 위한 Storage"}, {"name": "alertmanager", "version": "v0.23.0", "description": "Alert 관리를 위한 backend 서비스"}, {"name": "grafana", "version": "v6.50.7", "description": "모니터링 통합 포탈"}]}]' );
insert into stack_templates ( id, organization_id, name, description, version, cloud_service, platform, template, kube_version, kube_type, created_at, updated_at, services )
values ( '3696cb38-4da0-4235-97eb-b6eb15962bd1', 'master', 'AWS Standard (arm)', 'included LMA, SERVICE_MESH', 'v2', 'AWS', 'arm', 'aws-mar-msa-reference', 'v1.24', 'EKS', now(), now(), '[{"name": "Logging,Monitoring,Alerting", "type": "LMA", "applications": [{"name": "prometheus-stack", "version": "v.44.3.1", "description": "통계데이터 제공을 위한 backend  플랫폼"}, {"name": "elastic-system", "version": "v1.8.0", "description": "로그 데이터 적재를 위한 Storage"}, {"name": "alertmanager", "version": "v0.23.0", "description": "Alert 관리를 위한 backend 서비스"}, {"name": "grafana", "version": "v6.50.7", "description": "모니터링 통합 포탈"}]}, {"name": "MSA", "type": "SERVICE_MESH", "applications": [{"name": "istio", "version": "v1.13.1", "description": "MSA 플랫폼"}, {"name": "jagger", "version": "v2.27.1", "description": "분산 서비스간 트랜잭션 추적을 위한 로깅 플랫폼"}, {"name": "kiali", "version": "v1.45.1", "description": "MSA 통합 모니터링포탈"}]}]' );
