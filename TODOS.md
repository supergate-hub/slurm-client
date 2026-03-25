# TODOS

## P2: MCP 고도화 2단계 — GPU allocation + Queue depth
slurmrestd GRES 엔드포인트 매핑 조사 후 `slurm_gpu_allocation`, `slurm_queue_depth` MCP 툴 구현. 실제 GPU utilization이 아닌 GRES 할당 현황(allocated/total) 기반.
- Effort: M (human) → S (CC)
- Depends on: Multi-cluster plan 완료
- Context: Design doc의 "고도화 2단계" 참조. api/v0040의 node GRES 필드 확인 필요.

## P2: RBAC / Audit logging for MCP Server
AI 에이전트가 HPC를 제어하려면 operation scoping(read-only vs full-access) + audit trail 필요. SSE transport로 네트워크 노출 시 특히 중요. 프로덕션 HPC 환경의 보안/컴플라이언스 요구사항.
- Effort: L (human) → M (CC)
- Depends on: SSE transport + bearer token auth 완료
- Context: Outside voice review에서 지적. slurmrestd 자체의 토큰 권한과 별개로 MCP 레벨 접근 제어.

## P3: Prometheus/DCGM 연동 데이터 소스
실시간 GPU utilization(nvidia-smi 수준)은 slurmrestd가 제공하지 않음. Prometheus/DCGM exporter에서 데이터 수집하는 별도 어댑터 필요.
- Effort: L (human) → M (CC)
- Depends on: 고도화 2단계 완료
- Context: MCP 고도화 3단계로 위치. 별도 데이터 소스 연동은 아키텍처 확장 필요.
