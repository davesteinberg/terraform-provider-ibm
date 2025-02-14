// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-beta-go-sdk/vpcbetav1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIbmIsSourceShare() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmIsSourceShareRead,

		Schema: map[string]*schema.Schema{
			"share_replica": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The replica file share identifier.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the file share is created.",
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this share.",
			},
			"encryption": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of encryption used for this file share.",
			},
			"encryption_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The key used to encrypt this file share. The CRN of the [Key Protect Root Key](https://cloud.ibm.com/docs/key-protect?topic=key-protect-getting-started-tutorial) or [Hyper Protect Crypto Service Root Key](https://cloud.ibm.com/docs/hs-crypto?topic=hs-crypto-get-started) for this resource.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this share.",
			},
			"iops": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum input/output operation performance bandwidth per second for the file share.",
			},
			"latest_job": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The latest job associated with this file share.This property will be absent if no jobs have been created for this file share.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the file share job.The enumerated values for this property will expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the file share job on which the unexpected property value was encountered.* `cancelled`: This job has been cancelled.* `failed`: This job has failed.* `queued`: This job is queued.* `running`: This job is running.* `succeeded`: This job completed successfully.",
						},
						"status_reasons": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The reasons for the file share job status (if any).The enumerated reason code values for this property will expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected reason code was encountered.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A snake case string succinctly identifying the status reason.",
									},
									"message": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "An explanation of the status reason.",
									},
									"more_info": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about this status reason.",
									},
								},
							},
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the file share job.The enumerated values for this property will expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the file share job on which the unexpected property value was encountered.* `replication_failover`: This is a share replication failover job.* `replication_init`: This is a share replication is initialization job.* `replication_split`: This is a share replication split job.",
						},
					},
				},
			},
			"lifecycle_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the file share.",
			},
			"profile": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The globally unique name of the profile this file share uses.",
			},
			"replica_share": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The replica file share for this source file share.This property will be present when the `replication_role` is `source`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this file share.",
						},
						"deleted": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this file share.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this file share.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique user-defined name for this file share.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			"replication_cron_spec": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The cron specification for the file share replication schedule.This property will be present when the `replication_role` is `replica`.",
			},
			"replication_role": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The replication role of the file share.* `none`: This share is not participating in replication.* `replica`: This share is a replication target.* `source`: This share is a replication source.",
			},
			"replication_status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The replication status of the file share.* `active`: This share is actively participating in replication, and the replica's data is up-to-date with the replication schedule.* `failover_pending`: This share is performing a replication failover.* `initializing`: This share is initializing replication.* `none`: This share is not participating in replication.* `split_pending`: This share is performing a replication split.",
			},
			"replication_status_reasons": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The reasons for the current replication status (if any).The enumerated reason code values for this property will expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected reason code was encountered.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A snake case string succinctly identifying the status reason.",
						},
						"message": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An explanation of the status reason.",
						},
						"more_info": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Link to documentation about this status reason.",
						},
					},
				},
			},
			"resource_group": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the resource group for this file share.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of resource referenced.",
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The size of the file share rounded up to the next gigabyte.",
			},
			"source_share": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The source file share for this replica file share.This property will be present when the `replication_role` is `replica`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this file share.",
						},
						"deleted": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this file share.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this file share.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique user-defined name for this file share.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			"share_targets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Mount targets for the file share.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deleted": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this share target.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this share target.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this share target.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of resource referenced.",
						},
					},
				},
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the share.",
			},
			"zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The globally unique name of the zone this file share will reside in.",
			},
			isFileShareAccessTags: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of access management tags",
			},
			isFileShareTags: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of tags",
			},
		},
	}
}

func dataSourceIbmIsSourceShareRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1BetaAPI()
	if err != nil {
		return diag.FromErr(err)
	}

	replicaShareId := d.Get("share_replica").(string)

	getShareSourceOptions := &vpcbetav1.GetShareSourceOptions{
		ShareID: &replicaShareId,
	}

	share, response, err := vpcClient.GetShareSourceWithContext(context, getShareSourceOptions)
	if err != nil {
		if response != nil {
			if response.StatusCode == 404 {
				d.SetId("")
			}
			log.Printf("[DEBUG] GetShareWithContext failed %s\n%s", err, response)
			return nil
		}
		log.Printf("[DEBUG] GetShareWithContext failed %s\n", err)
		return diag.FromErr(fmt.Errorf("[DEBUG] GetShareWithContext failed %s\n", err))
	}

	d.SetId(*share.ID)
	if err = d.Set("created_at", share.CreatedAt.String()); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("crn", share.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}
	if err = d.Set("encryption", share.Encryption); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting encryption: %s", err))
	}

	if share.EncryptionKey != nil {
		err = d.Set("encryption_key", *share.EncryptionKey.CRN)
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting encryption_key %s", err))
		}
	}
	if err = d.Set("href", share.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("iops", share.Iops); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting iops: %s", err))
	}

	if share.LatestJob != nil {
		err = d.Set("latest_job", dataSourceShareFlattenLatestJob(*share.LatestJob))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting latest_job %s", err))
		}
	}

	if err = d.Set("lifecycle_state", share.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting lifecycle_state: %s", err))
	}
	if err = d.Set("name", share.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}

	if share.Profile != nil {
		err = d.Set("profile", *share.Profile.Name)
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting profile %s", err))
		}
	}

	if share.ReplicaShare != nil {
		err = d.Set("replica_share", dataSourceShareFlattenReplicaShare(*share.ReplicaShare))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting replica_share %s", err))
		}
	}
	if err = d.Set("replication_cron_spec", share.ReplicationCronSpec); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting replication_cron_spec: %s", err))
	}
	if err = d.Set("replication_role", share.ReplicationRole); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting replication_role: %s", err))
	}
	if err = d.Set("replication_status", share.ReplicationStatus); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting replication_status: %s", err))
	}

	if share.ReplicationStatusReasons != nil {
		err = d.Set("replication_status_reasons", dataSourceShareFlattenReplicationStatusReasons(share.ReplicationStatusReasons))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting replication_status_reasons %s", err))
		}
	}
	if share.ResourceGroup != nil {
		err = d.Set("resource_group", *share.ResourceGroup.ID)
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting resource_group %s", err))
		}
	}
	if err = d.Set("resource_type", share.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}
	if err = d.Set("size", share.Size); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting size: %s", err))
	}
	if share.SourceShare != nil {
		err = d.Set("source_share", dataSourceShareFlattenSourceShare(*share.SourceShare))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting source_share %s", err))
		}
	}
	if share.MountTargets != nil {
		err = d.Set("share_targets", dataSourceShareFlattenTargets(share.MountTargets))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting targets %s", err))
		}
	}

	if share.Zone != nil {
		err = d.Set("zone", *share.Zone.Name)
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting zone %s", err))
		}
	}
	tags, err := flex.GetGlobalTagsUsingCRN(meta, *share.CRN, "", isUserTagType)
	if err != nil {
		log.Printf(
			"Error getting shares (%s) tags: %s", d.Id(), err)
	}

	accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *share.CRN, "", isAccessTagType)
	if err != nil {
		log.Printf(
			"Error getting shares (%s) access tags: %s", d.Id(), err)
	}

	d.Set(isFileShareTags, tags)
	d.Set(isFileShareAccessTags, accesstags)

	return nil
}
