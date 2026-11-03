package com.hpe.terraform.skeleton;

import com.morpheusdata.core.Plugin;

/**
 * Inert Morpheus plugin used by the HPE Terraform provider acceptance tests.
 *
 * It registers no providers and touches nothing on the appliance: uploading it
 * only exercises the plugin install / update / delete lifecycle that the
 * hpe_morpheus_plugin resource manages. The version is taken from the jar
 * manifest (Plugin-Version), so the same class is packaged twice to test the
 * update path.
 */
public class SkeletonPlugin extends Plugin {

    public static final String CODE = "hpe-terraform-skeleton-plugin";

    @Override
    public String getCode() {
        return CODE;
    }

    @Override
    public void initialize() {
        this.setName("HPE Terraform Skeleton");
        this.setDescription("Inert plugin used by the HPE Terraform provider acceptance tests");
        this.setAuthor("HPE");
    }

    @Override
    public void onDestroy() {
        // Nothing to clean up: the plugin registers no providers.
    }
}
